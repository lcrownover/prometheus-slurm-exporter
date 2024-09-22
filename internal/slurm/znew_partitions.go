//go:build 2311

package slurm

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

func PartitionsData() []byte {
	cmd := exec.Command("sinfo", "-h", "-o%R,%C")
	// allocated/idle/other/total
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out, _ := ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return out
}

func PartitionsPendingJobsData() []byte {
	cmd := exec.Command("squeue", "-a", "-r", "-h", "-o%P", "--states=PENDING")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out, _ := ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return out
}

func NewPartitionsMetrics() *partitionMetrics {
	return &partitionMetrics{0, 0, 0, 0, 0}
}

type partitionMetrics struct {
	cpus_allocated float64
	cpus_idle      float64
	cpus_other     float64
	cpus_total     float64
	jobs_pending   float64
}

// ParsePartitionsMetrics returns a map where the keys are the partition names and the values are a partitionMetrics struct
func ParsePartitionsMetrics(partitionResp types.V0040OpenapiPartitionResp, jobsResp types.V0040OpenapiJobInfoResp, nodesResp types.V0040OpenapiNodesResp) (map[string]*partitionMetrics, error) {
	partitions := make(map[string]*partitionMetrics)
	nodePartitions := make(map[string][]string)

	// first, store all the nodes and their partitions
	for _, n := range nodesResp.Nodes {
		nodeName, err := GetNodeName(n)
		if err != nil {
			return nil, fmt.Errorf("failed to get node name for partition metrics: %v", err)
		}
		nodePartitions[*nodeName] = GetNodePartitions(n)
	}

	// scan through partition data to get total cpus
	for _, p := range partitionResp.Partitions {
		partition, err := GetPartitionName(p)
		if err != nil {
			return nil, fmt.Errorf("failed to get partition name for partition metrics: %v", err)
		}
		_, exists := partitions[*partition]
		if !exists {
			partitions[*partition] = NewPartitionsMetrics()
		}

		// cpu total
		total, err := GetPartitionTotalCPUs(p)
		if err != nil {
			return nil, fmt.Errorf("failed to collect cpu total for partition metrics: %v", err)
		}
		partitions[*partition].cpus_total = *total
	}

	// to get used and available cpus, we need to scan through the job list and categorize
	// each job by its partition, adding the cpus as we go
	for _, n := range nodesResp.Nodes {
		alloc_cpus := GetNodeAllocCPUs(n)
		idle_cpus := GetNodeIdleCPUs(n)
		nodePartitions := GetNodePartitions(n)
		for _, pname := range nodePartitions {
			partitions[pname].cpus_allocated += float64(alloc_cpus)
			partitions[pname].cpus_idle += float64(idle_cpus)
		}
	}

	// derive the other stat
	for i, p := range partitions {
		partitions[i].cpus_other = p.cpus_total - p.cpus_allocated - p.cpus_idle
	}

	// lastly, we need to get a count of pending jobs for the partition
	for _, j := range jobsResp.Jobs {
		pname, err := GetJobPartitionName(j)
		if err != nil {
			return nil, fmt.Errorf("failed to get job partition name for partition metrics: %v", err)
		}
		// partition name can be comma-separated, so we iterate through it
		pnames := strings.Split(*pname, ",")
		for _, pname := range pnames {
			slog.Info("job partition name", "name", pname)
			partitions[pname].jobs_pending += 1
		}
	}

	return partitions, nil
}

type PartitionsCollector struct {
	ctx       context.Context
	allocated *prometheus.Desc
	idle      *prometheus.Desc
	other     *prometheus.Desc
	pending   *prometheus.Desc
	total     *prometheus.Desc
}

func NewPartitionsCollector(ctx context.Context) *PartitionsCollector {
	labels := []string{"partition"}
	return &PartitionsCollector{
		ctx:       ctx,
		allocated: prometheus.NewDesc("slurm_partition_cpus_allocated", "Allocated CPUs for partition", labels, nil),
		idle:      prometheus.NewDesc("slurm_partition_cpus_idle", "Idle CPUs for partition", labels, nil),
		other:     prometheus.NewDesc("slurm_partition_cpus_other", "Other CPUs for partition", labels, nil),
		pending:   prometheus.NewDesc("slurm_partition_jobs_pending", "Pending jobs for partition", labels, nil),
		total:     prometheus.NewDesc("slurm_partition_cpus_total", "Total CPUs for partition", labels, nil),
	}
}

func (pc *PartitionsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- pc.allocated
	ch <- pc.idle
	ch <- pc.other
	ch <- pc.pending
	ch <- pc.total
}

func (pc *PartitionsCollector) Collect(ch chan<- prometheus.Metric) {
	partitionRespBytes, err := api.GetSlurmRestPartitionsResponse(pc.ctx)
	if err != nil {
		slog.Error("failed to get partitions response for partitions metrics", "error", err)
		return
	}
	partitionsResp, err := api.UnmarshalPartitionsResponse(partitionRespBytes)
	if err != nil {
		slog.Error("failed to unmarshal partitions response for partitions metrics", "error", err)
		return
	}
	jobRespBytes, err := api.GetSlurmRestJobsResponse(pc.ctx)
	if err != nil {
		slog.Error("failed to get jobs response for partitions metrics", "error", err)
		return
	}
	jobsResp, err := api.UnmarshalJobsResponse(jobRespBytes)
	if err != nil {
		slog.Error("failed to unmarshal jobs response for partitions metrics", "error", err)
		return
	}
	nodeRespBytes, err := api.GetSlurmRestNodesResponse(pc.ctx)
	if err != nil {
		slog.Error("failed to get nodes response for partition metrics", "error", err)
		return
	}
	nodesResp, err := api.UnmarshalNodesResponse(nodeRespBytes)
	if err != nil {
		slog.Error("failed to unmarshal nodes response for partition metrics", "error", err)
		return
	}
	pm, err := ParsePartitionsMetrics(*partitionsResp, *jobsResp, *nodesResp)
	if err != nil {
		slog.Error("failed to collect partitions metrics", "error", err)
		return
	}
	for p := range pm {
		if pm[p].cpus_allocated > 0 {
			ch <- prometheus.MustNewConstMetric(pc.allocated, prometheus.GaugeValue, pm[p].cpus_allocated, p)
		}
		if pm[p].cpus_idle > 0 {
			ch <- prometheus.MustNewConstMetric(pc.idle, prometheus.GaugeValue, pm[p].cpus_idle, p)
		}
		if pm[p].cpus_other > 0 {
			ch <- prometheus.MustNewConstMetric(pc.other, prometheus.GaugeValue, pm[p].cpus_other, p)
		}
		if pm[p].cpus_total > 0 {
			ch <- prometheus.MustNewConstMetric(pc.total, prometheus.GaugeValue, pm[p].cpus_total, p)
		}
		if pm[p].jobs_pending > 0 {
			ch <- prometheus.MustNewConstMetric(pc.pending, prometheus.GaugeValue, pm[p].jobs_pending, p)
		}
	}
}
