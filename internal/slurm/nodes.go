package slurm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

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
func ParseNodesMetrics(nodesResp types.V0040OpenapiNodesResp) (*nodesMetrics, error) {
	nm := NewNodesMetrics()

	for _, n := range nodesResp.Nodes {
		nodeStates, err := GetNodeStates(n)
		if err != nil {
			return nil, fmt.Errorf("failed to get node state for nodes metrics: %v", err)
		}

		for _, ns := range *nodeStates {
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

/*
 * Implement the Prometheus Collector interface and feed the
 * Slurm scheduler metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */

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
	nm, err := ParseNodesMetrics(*nodesResp)
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
