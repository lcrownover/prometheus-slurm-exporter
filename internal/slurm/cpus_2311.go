//go:build 2311

package slurm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/akyoto/cache"
	openapi "github.com/lcrownover/openapi-slurm-23-11"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

// CPU metrics collector
type CPUsCollector struct {
	ctx   context.Context
	alloc *prometheus.Desc
	idle  *prometheus.Desc
	other *prometheus.Desc
	total *prometheus.Desc
}

// NewCPUsCollector creates a new CPUsCollector
func NewCPUsCollector(ctx context.Context) *CPUsCollector {
	return &CPUsCollector{
		ctx:   ctx,
		alloc: prometheus.NewDesc("slurm_cpus_alloc", "Allocated CPUs", nil, nil),
		idle:  prometheus.NewDesc("slurm_cpus_idle", "Idle CPUs", nil, nil),
		other: prometheus.NewDesc("slurm_cpus_other", "Mix CPUs", nil, nil),
		total: prometheus.NewDesc("slurm_cpus_total", "Total CPUs", nil, nil),
	}
}

func (cc *CPUsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.other
	ch <- cc.total
}

func (cc *CPUsCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := cc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	jobsRespBytes, found := apiCache.Get("jobs")
	if !found {
		slog.Error("failed to get jobs response for users metrics from cache")
		return
	}
	jobsResp, err := api.UnmarshalJobsResponse(jobsRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal jobs response for cpu metrics", "error", err)
		return
	}
	nodeRespBytes, found := apiCache.Get("nodes")
	if !found {
		slog.Error("failed to get nodes response for cpu metrics from cache")
		return
	}
	nodesResp, err := api.UnmarshalNodesResponse(nodeRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal nodes response for cpu metrics", "error", err)
		return
	}
	cm, err := ParseCPUsMetrics(*nodesResp, *jobsResp)
	if err != nil {
		slog.Error("failed to collect cpus metrics", "error", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
}

type cpusMetrics struct {
	alloc float64
	idle  float64
	other float64
	total float64
}

func NewCPUsMetrics() *cpusMetrics {
	return &cpusMetrics{}
}

// ParseCPUMetrics pulls out total cluster cpu states of alloc,idle,other,total
func ParseCPUsMetrics(nodesResp openapi.V0040OpenapiNodesResp, jobsResp openapi.V0040OpenapiJobInfoResp) (*cpusMetrics, error) {
	cm := NewCPUsMetrics()
	for _, j := range jobsResp.Jobs {
		state, err := GetJobState(j)
		if err != nil {
			slog.Error("failed to get job state", "error", err)
			continue
		}
		cpus, err := GetJobCPUs(j)
		if err != nil {
			slog.Error("failed to get job cpus", "error", err)
			continue
		}
		// alloc is easy, we just add up all the cpus in the "Running" job state
		if *state == types.JobStateRunning {
			cm.alloc += *cpus
		}
	}
	// total is just the total number of cpus in the cluster
	nodes := nodesResp.Nodes
	for _, n := range nodes {
		if *n.Cpus == 1 {
			// TODO: This probably needs to be a call to partitions to get all nodes
			// in a partition, then add the nodes CPU values up for this field.
			// In our environment, nodes that exist (need slurm commands) get
			// put into slurm without being assigned a partition, but slurm
			// seems to track these systems with cpus=1.
			// This isn't a problem unless your site has nodes with a single CPU.
			continue
		}
		cpus := float64(*n.Cpus)
		cm.total += cpus

		nodeStates, err := GetNodeStates(n)
		if err != nil {
			return nil, fmt.Errorf("failed to get node state for cpu metrics: %v", err)
		}
		for _, ns := range *nodeStates {
			if ns == types.NodeStateMix || ns == types.NodeStateAlloc || ns == types.NodeStateIdle {
				// TODO: This calculate is scuffed. In our 17k core environment, it's
				// reporting ~400 more than the `sinfo -h -o '%C'` command.
				// Gotta figure this one out.
				idle_cpus := float64(*n.AllocIdleCpus)
				cm.idle += idle_cpus
			}
		}
	}
	// Assumedly, this should be fine.
	cm.other = cm.total - cm.idle - cm.alloc
	return cm, nil
}
