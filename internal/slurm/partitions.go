package slurm

import (
	"context"
	"log/slog"
	"strings"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

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
	apiCache := pc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	partitionsRespBytes, found := apiCache.Get("partitions")
	if !found {
		slog.Error("failed to get partitions response for partitions metrics from cache")
		return
	}
	jobsRespBytes, found := apiCache.Get("jobs")
	if !found {
		slog.Error("failed to get jobs response for users metrics from cache")
		return
	}
	nodesRespBytes, found := apiCache.Get("nodes")
	if !found {
		slog.Error("failed to get nodes response for cpu metrics from cache")
		return
	}
	partitionsData, err := api.ProcessPartitionsResponse(partitionsRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to process partitions data for partitions metrics", "error", err)
		return
	}
	jobsData, err := api.ProcessJobsResponse(jobsRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to process jobs data for partitions metrics", "error", err)
		return
	}
	nodesData, err := api.ProcessNodesResponse(nodesRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to process nodes data for partitions metrics", "error", err)
		return
	}
	pm, err := ParsePartitionsMetrics(partitionsData, jobsData, nodesData)
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
func ParsePartitionsMetrics(partitionsData *api.PartitionsData, jobsData *api.JobsData, nodesData *api.NodesData) (map[string]*partitionMetrics, error) {
	partitions := make(map[string]*partitionMetrics)
	nodePartitions := make(map[string][]string)

	// first, scan through partition data to easily get total cpus
	for _, p := range partitionsData.Partitions {
		_, exists := partitions[p.Name]
		if !exists {
			partitions[p.Name] = NewPartitionsMetrics()
		}

		// cpu total
		partitions[p.Name].cpus_total = float64(p.Cpus)
	}

	// we need to gather cpus from the nodes perspective because a node can
	// be a member of multiple partitions, running a job in one partition, and
	// we want to see that there are allocated cpus on the other partition because
	// of the shared node.
	for _, n := range nodesData.Nodes {
		nodePartitions[n.Name] = n.Partitions
	}

	// to get used and available cpus, we need to scan through the job list and categorize
	// each job by its partition, adding the cpus as we go
	for _, n := range nodesData.Nodes {
		alloc_cpus := n.AllocCpus
		idle_cpus := n.AllocIdleCpus
		nodePartitionNames := n.Partitions
		for _, partitionName := range nodePartitionNames {
			// this needs to exist to handle the test data provided by SLURM
			// where the nodes response example data does not correspond to
			// the partitions response example data. in real data, the
			// partition names should already exist in the map
			_, exists := partitions[partitionName]
			if !exists {
				partitions[partitionName] = NewPartitionsMetrics()
			}

			partitions[partitionName].cpus_allocated += float64(alloc_cpus)
			partitions[partitionName].cpus_idle += float64(idle_cpus)
		}
	}

	// derive the other stat
	for i, p := range partitions {
		partitions[i].cpus_other = p.cpus_total - p.cpus_allocated - p.cpus_idle
	}

	// lastly, we need to get a count of pending jobs for the partition
	for _, j := range jobsData.Jobs {
		// partition name can be comma-separated, so we iterate through it
		pnames := strings.Split(j.Partition, ",")
		for _, partitionName := range pnames {
			// this needs to exist to handle the test data provided by SLURM
			// where the nodes response example data does not correspond to
			// the partitions response example data. in real data, the
			// partition names should already exist in the map
			_, exists := partitions[partitionName]
			if !exists {
				partitions[partitionName] = NewPartitionsMetrics()
			}
			partitions[partitionName].jobs_pending += 1
		}
	}

	return partitions, nil
}
