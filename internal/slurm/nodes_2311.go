//go:build 2311

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
