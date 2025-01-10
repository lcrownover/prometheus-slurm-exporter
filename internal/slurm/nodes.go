package slurm

import (
	"context"
	"log/slog"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

type NodesCollector struct {
	ctx   context.Context
	alloc *prometheus.Desc
	comp  *prometheus.Desc
	down  *prometheus.Desc
	drain *prometheus.Desc
	err   *prometheus.Desc
	fail  *prometheus.Desc
	idle  *prometheus.Desc
	maint *prometheus.Desc
	mix   *prometheus.Desc
	resv  *prometheus.Desc
}

func NewNodesCollector(ctx context.Context) *NodesCollector {
	return &NodesCollector{
		ctx:   ctx,
		alloc: prometheus.NewDesc("slurm_nodes_alloc", "Allocated nodes", nil, nil),
		comp:  prometheus.NewDesc("slurm_nodes_comp", "Completing nodes", nil, nil),
		down:  prometheus.NewDesc("slurm_nodes_down", "Down nodes", nil, nil),
		drain: prometheus.NewDesc("slurm_nodes_drain", "Drain nodes", nil, nil),
		err:   prometheus.NewDesc("slurm_nodes_err", "Error nodes", nil, nil),
		fail:  prometheus.NewDesc("slurm_nodes_fail", "Fail nodes", nil, nil),
		idle:  prometheus.NewDesc("slurm_nodes_idle", "Idle nodes", nil, nil),
		maint: prometheus.NewDesc("slurm_nodes_maint", "Maint nodes", nil, nil),
		mix:   prometheus.NewDesc("slurm_nodes_mix", "Mix nodes", nil, nil),
		resv:  prometheus.NewDesc("slurm_nodes_resv", "Reserved nodes", nil, nil),
	}
}

func (nc *NodesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.alloc
	ch <- nc.comp
	ch <- nc.down
	ch <- nc.drain
	ch <- nc.err
	ch <- nc.fail
	ch <- nc.idle
	ch <- nc.maint
	ch <- nc.mix
	ch <- nc.resv
}

func (nc *NodesCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := nc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	nodesRespBytes, found := apiCache.Get("nodes")
	if !found {
		slog.Error("failed to get nodes response for cpu metrics from cache")
		return
	}
	nodesData, err := api.ProcessNodesResponse(nodesRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to process nodes response for nodes metrics", "error", err)
		return
	}
	nm, err := ParseNodesMetrics(nodesData)
	if err != nil {
		slog.Error("failed to collect nodes metrics", "error", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(nc.alloc, prometheus.GaugeValue, nm.alloc)
	ch <- prometheus.MustNewConstMetric(nc.comp, prometheus.GaugeValue, nm.comp)
	ch <- prometheus.MustNewConstMetric(nc.down, prometheus.GaugeValue, nm.down)
	ch <- prometheus.MustNewConstMetric(nc.drain, prometheus.GaugeValue, nm.drain)
	ch <- prometheus.MustNewConstMetric(nc.err, prometheus.GaugeValue, nm.err)
	ch <- prometheus.MustNewConstMetric(nc.fail, prometheus.GaugeValue, nm.fail)
	ch <- prometheus.MustNewConstMetric(nc.idle, prometheus.GaugeValue, nm.idle)
	ch <- prometheus.MustNewConstMetric(nc.maint, prometheus.GaugeValue, nm.maint)
	ch <- prometheus.MustNewConstMetric(nc.mix, prometheus.GaugeValue, nm.mix)
	ch <- prometheus.MustNewConstMetric(nc.resv, prometheus.GaugeValue, nm.resv)
}

type nodesMetrics struct {
	alloc float64
	comp  float64
	down  float64
	drain float64
	err   float64
	fail  float64
	idle  float64
	maint float64
	mix   float64
	resv  float64
}

func NewNodesMetrics() *nodesMetrics {
	return &nodesMetrics{}
}

// ParseNodesMetrics iterates through node response objects and tallies up
// nodes based on their state
func ParseNodesMetrics(nodesData *api.NodesData) (*nodesMetrics, error) {
	nm := NewNodesMetrics()

	for _, n := range nodesData.Nodes {
		for _, ns := range n.States {
			switch ns {
			case types.NodeStateAlloc:
				nm.alloc += 1
			case types.NodeStateComp:
				nm.comp += 1
			case types.NodeStateDown:
				nm.down += 1
			case types.NodeStateDrain:
				nm.drain += 1
			case types.NodeStateErr:
				nm.err += 1
			case types.NodeStateFail:
				nm.fail += 1
			case types.NodeStateIdle:
				nm.idle += 1
			case types.NodeStateMaint:
				nm.maint += 1
			case types.NodeStateMix:
				nm.mix += 1
			case types.NodeStateResv:
				nm.resv += 1
			}
		}
	}

	return nm, nil
}
