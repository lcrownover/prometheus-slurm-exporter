package slurm

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

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
func ParseNodeMetrics(nodesResp types.V0040OpenapiNodesResp) (map[string]*nodeMetrics, error) {
	nodeMap := make(map[string]*nodeMetrics)

	for _, n := range nodesResp.Nodes {
		nodeName := *n.Hostname
		nodeMap[nodeName] = &nodeMetrics{0, 0, 0, 0, 0, 0, ""}

		// state
		nodeStatesStr, err := GetNodeStatesString(n, "|")
		if err != nil {
			return nil, fmt.Errorf("failed to get node state: %v", err)
		}
		nodeMap[nodeName].nodeStatus = nodeStatesStr

		// memory
		nodeMap[nodeName].memAlloc = GetNodeAllocMemory(n)
		nodeMap[nodeName].memTotal = GetNodeTotalMemory(n)

		// cpu
		nodeMap[nodeName].cpuAlloc = GetNodeAllocCPUs(n)
		nodeMap[nodeName].cpuIdle = GetNodeIdleCPUs(n)
		nodeMap[nodeName].cpuOther = GetNodeOtherCPUs(n)
		nodeMap[nodeName].cpuTotal = GetNodeTotalCPUs(n)
	}

	return nodeMap, nil
}

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
	nodeRespBytes, err := api.GetSlurmRestNodesResponse(nc.ctx)
	if err != nil {
		slog.Error("failed to get nodes response for cpu metrics", "error", err)
		return
	}
	nodesResp, err := api.UnmarshalNodesResponse(nodeRespBytes)
	if err != nil {
		slog.Error("failed to unmarshal nodes response for cpu metrics", "error", err)
		return
	}
	nm, err := ParseNodeMetrics(*nodesResp)
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

/*

OLD CODE HERE FOR NOW

*/

// NodeMetrics stores metrics for each node
type NodeMetrics struct {
	memAlloc   uint64
	memTotal   uint64
	cpuAlloc   uint64
	cpuIdle    uint64
	cpuOther   uint64
	cpuTotal   uint64
	nodeStatus string
}

func NodeGetMetricsOld() map[string]*NodeMetrics {
	return ParseNodeMetricsOld(NodeDataOld())
}

// ParseNodeMetricsOld takes the output of sinfo with node data
// It returns a map of metrics per node
func ParseNodeMetricsOld(input []byte) map[string]*NodeMetrics {
	nodes := make(map[string]*NodeMetrics)
	lines := strings.Split(string(input), "\n")

	// Sort and remove all the duplicates from the 'sinfo' output
	sort.Strings(lines)
	linesUniq := OldRemoveDuplicates(lines)

	for _, line := range linesUniq {
		node := strings.Fields(line)
		nodeName := node[0]
		nodeStatus := node[4] // mixed, allocated, etc.

		nodes[nodeName] = &NodeMetrics{0, 0, 0, 0, 0, 0, ""}

		memAlloc, _ := strconv.ParseUint(node[1], 10, 64)
		memTotal, _ := strconv.ParseUint(node[2], 10, 64)

		cpuInfo := strings.Split(node[3], "/")
		cpuAlloc, _ := strconv.ParseUint(cpuInfo[0], 10, 64)
		cpuIdle, _ := strconv.ParseUint(cpuInfo[1], 10, 64)
		cpuOther, _ := strconv.ParseUint(cpuInfo[2], 10, 64)
		cpuTotal, _ := strconv.ParseUint(cpuInfo[3], 10, 64)

		nodes[nodeName].memAlloc = memAlloc
		nodes[nodeName].memTotal = memTotal
		nodes[nodeName].cpuAlloc = cpuAlloc
		nodes[nodeName].cpuIdle = cpuIdle
		nodes[nodeName].cpuOther = cpuOther
		nodes[nodeName].cpuTotal = cpuTotal
		nodes[nodeName].nodeStatus = nodeStatus
	}

	return nodes
}

// NodeDataOld executes the sinfo command to get data for each node
// It returns the output of the sinfo command
func NodeDataOld() []byte {
	cmd := exec.Command("sinfo", "-h", "-N", "-O", "NodeList: ,AllocMem: ,Memory: ,CPUsState: ,StateLong:")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

type NodeCollectorOld struct {
	cpuAlloc *prometheus.Desc
	cpuIdle  *prometheus.Desc
	cpuOther *prometheus.Desc
	cpuTotal *prometheus.Desc
	memAlloc *prometheus.Desc
	memTotal *prometheus.Desc
}

// NewNodeCollectorOld creates a Prometheus collector to keep all our stats in
// It returns a set of collections for consumption
func NewNodeCollectorOld() *NodeCollectorOld {
	labels := []string{"node", "status"}

	return &NodeCollectorOld{
		cpuAlloc: prometheus.NewDesc("slurm_old_node_cpu_alloc", "Allocated CPUs per node", labels, nil),
		cpuIdle:  prometheus.NewDesc("slurm_old_node_cpu_idle", "Idle CPUs per node", labels, nil),
		cpuOther: prometheus.NewDesc("slurm_old_node_cpu_other", "Other CPUs per node", labels, nil),
		cpuTotal: prometheus.NewDesc("slurm_old_node_cpu_total", "Total CPUs per node", labels, nil),
		memAlloc: prometheus.NewDesc("slurm_old_node_mem_alloc", "Allocated memory per node", labels, nil),
		memTotal: prometheus.NewDesc("slurm_old_node_mem_total", "Total memory per node", labels, nil),
	}
}

// Send all metric descriptions
func (nc *NodeCollectorOld) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.cpuAlloc
	ch <- nc.cpuIdle
	ch <- nc.cpuOther
	ch <- nc.cpuTotal
	ch <- nc.memAlloc
	ch <- nc.memTotal
}

func (nc *NodeCollectorOld) Collect(ch chan<- prometheus.Metric) {
	nodes := NodeGetMetricsOld()
	for node := range nodes {
		ch <- prometheus.MustNewConstMetric(nc.cpuAlloc, prometheus.GaugeValue, float64(nodes[node].cpuAlloc), node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuIdle, prometheus.GaugeValue, float64(nodes[node].cpuIdle), node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuOther, prometheus.GaugeValue, float64(nodes[node].cpuOther), node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuTotal, prometheus.GaugeValue, float64(nodes[node].cpuTotal), node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.memAlloc, prometheus.GaugeValue, float64(nodes[node].memAlloc), node, nodes[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.memTotal, prometheus.GaugeValue, float64(nodes[node].memTotal), node, nodes[node].nodeStatus)
	}
}
