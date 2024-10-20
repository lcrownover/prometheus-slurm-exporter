package slurm

import (
	"context"
	"log/slog"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

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
	apiCache := fsc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	sharesRespBytes, found := apiCache.Get("shares")
	if !found {
		slog.Error("failed to get shares response for fair share metrics from cache")
		return
	}

	sharesData, err := api.ExtractSharesData(sharesRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to extract shares response for fair share metrics", "error", err)
		return
	}
	fsm, err := ParseFairShareMetrics(sharesData)
	if err != nil {
		slog.Error("failed to collect fair share metrics", "error", err)
		return
	}
	for f := range fsm {
		ch <- prometheus.MustNewConstMetric(fsc.fairshare, prometheus.GaugeValue, fsm[f].fairshare, f)
	}
}

type fairShareMetrics struct {
	fairshare float64
}

func NewFairShareMetrics() *fairShareMetrics {
	return &fairShareMetrics{}
}

func ParseFairShareMetrics(sharesData *api.SharesData) (map[string]*fairShareMetrics, error) {
	accounts := make(map[string]*fairShareMetrics)
	for _, s := range sharesData.Shares {
		account := s.Name
		if account == "root" {
			// we don't care about the root account
			continue
		}
		if _, exists := accounts[account]; !exists {
			accounts[account] = NewFairShareMetrics()
		}
		accounts[account].fairshare = s.EffectiveUsage
	}
	return accounts, nil
}
