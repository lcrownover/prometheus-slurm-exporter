package slurm

import (
	"context"
	"log/slog"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

type fairShareMetrics struct {
	fairshare float64
}

func NewFairShareMetrics() *fairShareMetrics {
	return &fairShareMetrics{}
}

func ParseFairShareMetrics(sharesResp types.V0040OpenapiSharesResp) (map[string]*fairShareMetrics, error) {
	accounts := make(map[string]*fairShareMetrics)
	for _, s := range *sharesResp.Shares.Shares {
		account := *s.Name
		if _, exists := accounts[account]; !exists {
			accounts[account] = NewFairShareMetrics()
		}
		// TODO: check if the level is the right value here,
		// there might be some other property that matches the
		// previous value from the old share info code
		accounts[account].fairshare = *s.Fairshare.Level
	}
	return accounts, nil
}

type FairShareCollector struct {
	ctx       context.Context
	fairshare *prometheus.Desc
}

func NewFairShareCollector(ctx context.Context) *FairShareCollector {
	labels := []string{"account"}
	return &FairShareCollector{
		ctx:       ctx,
		fairshare: prometheus.NewDesc("slurm_account_fairshare", "FairShare for account", labels, nil),
	}
}

func (fsc *FairShareCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- fsc.fairshare
}

func (fsc *FairShareCollector) Collect(ch chan<- prometheus.Metric) {
	sharesRespBytes, err := api.GetSlurmRestSharesResponse(fsc.ctx)
	if err != nil {
		slog.Error("failed to get shares response for fair share metrics", "error", err)
		return
	}
	sharesResp, err := api.UnmarshalSharesResponse(sharesRespBytes)
	if err != nil {
		slog.Error("failed to unmarshal shares response for fair share metrics", "error", err)
		return
	}
	fsm, err := ParseFairShareMetrics(*sharesResp)
	if err != nil {
		slog.Error("failed to collect fair share metrics", "error", err)
		return
	}
	for f := range fsm {
		ch <- prometheus.MustNewConstMetric(fsc.fairshare, prometheus.GaugeValue, fsm[f].fairshare, f)
	}
}
