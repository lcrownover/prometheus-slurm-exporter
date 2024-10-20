package slurm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

type NodeCollector struct {
	ctx      context.Context
	cpuAlloc *prometheus.Desc
	cpuIdle  *prometheus.Desc
	cpuOther *prometheus.Desc
	cpuTotal *prometheus.Desc
	memAlloc *prometheus.Desc
	memTotal *prometheus.Desc
}

// NewNodeCollectorOld creates a Prometheus collector to keep all our stats in
// It returns a set of collections for consumption
func NewNodeCollector(ctx context.Context) *NodeCollector {
	labels := []string{"node", "status"}

	return &NodeCollector{
		ctx:      ctx,
		cpuAlloc: prometheus.NewDesc("slurm_node_cpu_alloc", "Allocated CPUs per node", labels, nil),
		cpuIdle:  prometheus.NewDesc("slurm_node_cpu_idle", "Idle CPUs per node", labels, nil),
		cpuOther: prometheus.NewDesc("slurm_node_cpu_other", "Other CPUs per node", labels, nil),
		cpuTotal: prometheus.NewDesc("slurm_node_cpu_total", "Total CPUs per node", labels, nil),
		memAlloc: prometheus.NewDesc("slurm_node_mem_alloc", "Allocated memory per node", labels, nil),
		memTotal: prometheus.NewDesc("slurm_node_mem_total", "Total memory per node", labels, nil),
	}
}

// Send all metric descriptions
func (nc *NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.cpuAlloc
	ch <- nc.cpuIdle
	ch <- nc.cpuOther
	ch <- nc.cpuTotal
	ch <- nc.memAlloc
	ch <- nc.memTotal
}

func (nc *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := nc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	nodesRespBytes, found := apiCache.Get("nodes")
	if !found {
		slog.Error("failed to get nodes response for cpu metrics from cache")
		return
	}
	nodesData, err := api.ExtractNodesData(nodesRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to extract nodes response for cpu metrics", "error", err)
		return
	}
	nm, err := ParseNodeMetrics(nodesData)
	if err != nil {
		slog.Error("failed to collect nodes metrics", "error", err)
		return
	}
	for node := range nm {
		ch <- prometheus.MustNewConstMetric(nc.cpuAlloc, prometheus.GaugeValue, float64(nm[node].cpuAlloc), node, nm[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuIdle, prometheus.GaugeValue, float64(nm[node].cpuIdle), node, nm[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuOther, prometheus.GaugeValue, float64(nm[node].cpuOther), node, nm[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuTotal, prometheus.GaugeValue, float64(nm[node].cpuTotal), node, nm[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.memAlloc, prometheus.GaugeValue, float64(nm[node].memAlloc), node, nm[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.memTotal, prometheus.GaugeValue, float64(nm[node].memTotal), node, nm[node].nodeStatus)
	}
}

// NodeMetrics stores metrics for each node
type nodeMetrics struct {
	memAlloc   uint64
	memTotal   uint64
	cpuAlloc   uint64
	cpuIdle    uint64
	cpuOther   uint64
	cpuTotal   uint64
	nodeStatus string
}

func NewNodeMetrics() *nodeMetrics {
	return &nodeMetrics{}
}

// ParseNodeMetrics takes the output of sinfo with node data
// It returns a map of metrics per node
func ParseNodeMetrics(nodesData *api.NodesData) (map[string]*nodeMetrics, error) {
	nodeMap := make(map[string]*nodeMetrics)

	for _, n := range nodesData.Nodes {
		nodeName := n.Hostname
		nodeMap[nodeName] = &nodeMetrics{0, 0, 0, 0, 0, 0, ""}

		// state
		nodeStatesStr, err := n.GetNodeStatesString("|")
		if err != nil {
			return nil, fmt.Errorf("failed to get node state: %v", err)
		}
		nodeMap[nodeName].nodeStatus = nodeStatesStr

		// memory
		nodeMap[nodeName].memAlloc = uint64(n.AllocMemory)
		nodeMap[nodeName].memTotal = uint64(n.RealMemory)

		// cpu
		nodeMap[nodeName].cpuAlloc = uint64(n.AllocCpus)
		nodeMap[nodeName].cpuIdle = uint64(n.AllocIdleCpus)
		nodeMap[nodeName].cpuOther = uint64(n.OtherCpus)
		nodeMap[nodeName].cpuTotal = uint64(n.Cpus)
	}

	return nodeMap, nil
}
