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

type NodeStatusCollector struct {
	ctx    context.Context
	status *prometheus.Desc
}

// NewNodeStatusCollectorOld creates a Prometheus collector to keep all our stats in
// It returns a set of collections for consumption
func NewNodeStatusCollector(ctx context.Context) *NodeStatusCollector {
	labels := []string{"node", "status"}

	return &NodeStatusCollector{
		ctx:    ctx,
		status: prometheus.NewDesc("slurm_nodestatus", "Status enum for node", labels, nil),
	}
}

// Send all metric descriptions
func (nc *NodeStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.status
}

func (nc *NodeStatusCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := nc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	nodesRespBytes, found := apiCache.Get("nodes")
	if !found {
		slog.Error("failed to get nodes response for cpu metrics from cache")
		return
	}
	nodesData, err := api.ProcessNodesResponse(nodesRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to process nodes response for node metrics", "error", err)
		return
	}
	nm, err := ParseNodeStatusMetrics(nodesData)
	if err != nil {
		slog.Error("failed to collect nodes metrics", "error", err)
		return
	}
	for node := range nm {
		// TODO: implement
		statusStr := getNodeStatusFromStatusInt(nm[node].status)
		ch <- prometheus.MustNewConstMetric(nc.status, prometheus.GaugeValue, float64(nm[node].status), node, statusStr)
	}
}

// NodeStatusMetrics stores metrics for each node
type nodeStatusMetrics struct {
	status uint64
}

func NewNodeStatusMetrics() *nodeStatusMetrics {
	return &nodeStatusMetrics{}
}

func ParseNodeStatusMetrics(nodesData *api.NodesData) (map[string]*nodeStatusMetrics, error) {
	nodeMap := make(map[string]*nodeStatusMetrics)

	for _, n := range nodesData.Nodes {
		nodeName := n.Hostname
		nodeMap[nodeName] = &nodeStatusMetrics{0}

		// state
		nodeStatesStr, err := n.GetNodeStatesString("|")
		if err != nil {
			return nil, fmt.Errorf("failed to get node state: %v", err)
		}
		// TODO: do something here to convert nodestatestr into an int representing the state
		nodeMap[nodeName].status = 0
	}

	return nodeMap, nil
}
