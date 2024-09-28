//go:build 2405

package slurm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/akyoto/cache"
	openapi "github.com/lcrownover/openapi-slurm-24-05"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

type GPUsCollector struct {
	ctx         context.Context
	alloc       *prometheus.Desc
	idle        *prometheus.Desc
	other       *prometheus.Desc
	total       *prometheus.Desc
	utilization *prometheus.Desc
}

func NewGPUsCollector(ctx context.Context) *GPUsCollector {
	return &GPUsCollector{
		ctx:         ctx,
		alloc:       prometheus.NewDesc("slurm_gpus_alloc", "Allocated GPUs", nil, nil),
		idle:        prometheus.NewDesc("slurm_gpus_idle", "Idle GPUs", nil, nil),
		other:       prometheus.NewDesc("slurm_gpus_other", "Other GPUs", nil, nil),
		total:       prometheus.NewDesc("slurm_gpus_total", "Total GPUs", nil, nil),
		utilization: prometheus.NewDesc("slurm_gpus_utilization", "Total GPU utilization", nil, nil),
	}
}

func (cc *GPUsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.other
	ch <- cc.total
	ch <- cc.utilization
}
func (cc *GPUsCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := cc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
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
	gm, err := ParseGPUsMetrics(*nodesResp)
	if err != nil {
		slog.Error("failed to collect gpus metrics", "error", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, gm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, gm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, gm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, gm.total)
	ch <- prometheus.MustNewConstMetric(cc.utilization, prometheus.GaugeValue, gm.utilization)
}

type gpusMetrics struct {
	alloc       float64
	idle        float64
	other       float64
	total       float64
	utilization float64
}

func NewGPUsMetrics() *gpusMetrics {
	return &gpusMetrics{}
}

// NOTES:
// node[gres] 		=> gpu:0 										# no gpus
// node[gres] 		=> gpu:nvidia_h100_80gb_hbm3:4(S:0-1) 			# 4 h100 gpus
// node[gres_used]  => gpu:nvidia_h100_80gb_hbm3:4(IDX:0-3) 		# 4 used gpus
// node[gres_used]  => gpu:nvidia_h100_80gb_hbm3:0(IDX:N/A) 		# 0 used gpus
// node[tres]		=> cpu=48,mem=1020522M,billing=48,gres/gpu=4	# 4 total gpus
// node[tres]		=> cpu=1,mem=1M,billing=1						# 0 total gpus
// node[tres_used]	=> cpu=48,mem=1020522M,billing=48,gres/gpu=4	# 4 used gpus
// node[tres_used]	=> cpu=1,mem=1M,billing=1						# 0 used gpus
//
// For tracking gpu resources, it looks like tres will be better. If I need to pull out per-gpu stats later,
// I'll have to use gres
//

// ParseGPUsMetrics iterates through node response objects and tallies up the total and
// allocated gpus, then derives idle and utilization from those numbers.
func ParseGPUsMetrics(nodesResp openapi.V0041OpenapiNodesResp) (*gpusMetrics, error) {
	gm := NewGPUsMetrics()
	nodes := nodesResp.Nodes
	for _, n := range nodes {
		totalGPUs, err := GetNodeGPUTotal(n)
		if err != nil {
			return nil, fmt.Errorf("failed to get total gpu count for node: %v", err)
		}
		allocGPUs, err := GetNodeGPUAllocated(n)
		if err != nil {
			return nil, fmt.Errorf("failed to get allocated gpu count for node: %v", err)
		}
		idleGPUs := totalGPUs - allocGPUs
		gm.total += float64(totalGPUs)
		gm.alloc += float64(allocGPUs)
		gm.idle += float64(idleGPUs)
	}
	// TODO: Do we really need an "other" field?
	// using TRES, it should be straightforward.
	gm.utilization = gm.alloc / gm.total
	return gm, nil
}
