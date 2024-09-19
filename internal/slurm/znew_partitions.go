//go:build 2311

package slurm

//
// import (
// 	"context"
// 	"io/ioutil"
// 	"log"
// 	"log/slog"
// 	"os/exec"
// 	"strconv"
// 	"strings"
//
// 	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
// 	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
// 	"github.com/prometheus/client_golang/prometheus"
// )
//
// func PartitionsData() []byte {
// 	cmd := exec.Command("sinfo", "-h", "-o%R,%C")
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := cmd.Start(); err != nil {
// 		log.Fatal(err)
// 	}
// 	out, _ := ioutil.ReadAll(stdout)
// 	if err := cmd.Wait(); err != nil {
// 		log.Fatal(err)
// 	}
// 	return out
// }
//
// func PartitionsPendingJobsData() []byte {
// 	cmd := exec.Command("squeue", "-a", "-r", "-h", "-o%P", "--states=PENDING")
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := cmd.Start(); err != nil {
// 		log.Fatal(err)
// 	}
// 	out, _ := ioutil.ReadAll(stdout)
// 	if err := cmd.Wait(); err != nil {
// 		log.Fatal(err)
// 	}
// 	return out
// }
//
// func NewPartitionsMetrics() *partitionMetrics {
// 	return &partitionMetrics{0, 0, 0, 0, 0}
// }
//
// type partitionMetrics struct {
// 	cpus_allocated float64
// 	cpus_idle      float64
// 	cpus_other     float64
// 	cpus_total     float64
// 	jobs_pending   float64
// }
//
// func ParsePartitionsMetrics(partitionResp types.V0040OpenapiPartitionResp, jobsResp types.V0040OpenapiJobInfoResp) (map[string]*partitionMetrics, error) {
// 	partitions := make(map[string]*partitionMetrics)
// 	for _, j := range jobsResp.Jobs {
// 		partition := *j.Partition
// 		_, exists := partitions[partition]
// 		if !exists {
// 			partitions[partition] = NewPartitionsMetrics()
// 		}
//
// 		// cpu
// 		allocated, _ := strconv.ParseFloat(strings.Split(states, "/")[0], 64)
// 		idle, _ := strconv.ParseFloat(strings.Split(states, "/")[1], 64)
// 		other, _ := strconv.ParseFloat(strings.Split(states, "/")[2], 64)
// 		total, _ := strconv.ParseFloat(strings.Split(states, "/")[3], 64)
// 		partitions[partition].cpus_allocated = allocated
// 		partitions[partition].cpus_idle = idle
// 		partitions[partition].cpus_other = other
// 		partitions[partition].cpus_total = total
// 	}
//
// 	// get list of pending jobs by partition name
// 	list := strings.Split(string(PartitionsPendingJobsData()), "\n")
// 	for _, partition := range list {
// 		// accumulate the number of pending jobs
// 		_, key := partitions[partition]
// 		if key {
// 			partitions[partition].jobs_pending += 1
// 		}
// 	}
//
// 	return partitions, nil
// }
//
// type PartitionsCollector struct {
// 	ctx       context.Context
// 	allocated *prometheus.Desc
// 	idle      *prometheus.Desc
// 	other     *prometheus.Desc
// 	pending   *prometheus.Desc
// 	total     *prometheus.Desc
// }
//
// func NewPartitionsCollector(ctx context.Context) *PartitionsCollector {
// 	labels := []string{"partition"}
// 	return &PartitionsCollector{
// 		ctx:       ctx,
// 		allocated: prometheus.NewDesc("slurm_partition_cpus_allocated", "Allocated CPUs for partition", labels, nil),
// 		idle:      prometheus.NewDesc("slurm_partition_cpus_idle", "Idle CPUs for partition", labels, nil),
// 		other:     prometheus.NewDesc("slurm_partition_cpus_other", "Other CPUs for partition", labels, nil),
// 		pending:   prometheus.NewDesc("slurm_partition_jobs_pending", "Pending jobs for partition", labels, nil),
// 		total:     prometheus.NewDesc("slurm_partition_cpus_total", "Total CPUs for partition", labels, nil),
// 	}
// }
//
// func (pc *PartitionsCollector) Describe(ch chan<- *prometheus.Desc) {
// 	ch <- pc.allocated
// 	ch <- pc.idle
// 	ch <- pc.other
// 	ch <- pc.pending
// 	ch <- pc.total
// }
//
// func (pc *PartitionsCollector) Collect(ch chan<- prometheus.Metric) {
// 	partitionRespBytes, err := api.GetSlurmRestPartitionsResponse(pc.ctx)
// 	if err != nil {
// 		slog.Error("failed to get partitions response for partitions metrics", "error", err)
// 		return
// 	}
// 	partitionsResp, err := api.UnmarshalPartitionsResponse(partitionRespBytes)
// 	if err != nil {
// 		slog.Error("failed to unmarshal partitions response for partitions metrics", "error", err)
// 		return
// 	}
// 	pm, err := ParsePartitionsMetrics(*partitionsResp)
// 	if err != nil {
// 		slog.Error("failed to collect partitions metrics", "error", err)
// 		return
// 	}
// 	for p := range pm {
// 		if pm[p].cpus_allocated > 0 {
// 			ch <- prometheus.MustNewConstMetric(pc.allocated, prometheus.GaugeValue, pm[p].allocated, p)
// 		}
// 		if pm[p].cpus_idle > 0 {
// 			ch <- prometheus.MustNewConstMetric(pc.idle, prometheus.GaugeValue, pm[p].idle, p)
// 		}
// 		if pm[p].cpus_other > 0 {
// 			ch <- prometheus.MustNewConstMetric(pc.other, prometheus.GaugeValue, pm[p].other, p)
// 		}
// 		if pm[p].cpus_total > 0 {
// 			ch <- prometheus.MustNewConstMetric(pc.total, prometheus.GaugeValue, pm[p].total, p)
// 		}
// 		if pm[p].jobs_pending > 0 {
// 			ch <- prometheus.MustNewConstMetric(pc.pending, prometheus.GaugeValue, pm[p].pending, p)
// 		}
// 	}
// }
