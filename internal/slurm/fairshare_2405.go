//go:build 2405

package slurm

import (
	"context"
	"log/slog"
	"strings"

	"github.com/akyoto/cache"
	openapi "github.com/lcrownover/openapi-slurm-24-05"
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
	// this is disgusting but the response has values of "Infinity" which are
	// not json unmarshal-able, so I manually replace all the "Infinity"s with the correct
	// float64 value that represents Infinity.
	// this will be fixed in v0.0.42
	// https://support.schedmd.com/show_bug.cgi?id=20817
	//
	// https://github.com/lcrownover/prometheus-slurm-exporter/issues/8
	// also reported that folks are getting "inf" back, so I'll protect for that too
	sharesRespString := string(sharesRespBytes.([]byte))
	maxFloatStr := "1.7976931348623157e+308"
	// replacing the longer strings first should prevent any partial replacements
	sharesRespString = strings.ReplaceAll(sharesRespString, "Infinity", maxFloatStr)
	sharesRespString = strings.ReplaceAll(sharesRespString, "infinity", maxFloatStr)
	// sometimes it'd return "inf", so let's cover for that too.
	sharesRespString = strings.ReplaceAll(sharesRespString, "Inf", maxFloatStr)
	sharesRespString = strings.ReplaceAll(sharesRespString, "inf", maxFloatStr)
	sharesRespBytes = []byte(sharesRespString)
	// end hack

	sharesResp, err := api.UnmarshalSharesResponse(sharesRespBytes.([]byte))
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

type fairShareMetrics struct {
	fairshare float64
}

func NewFairShareMetrics() *fairShareMetrics {
	return &fairShareMetrics{}
}

func ParseFairShareMetrics(sharesResp openapi.SlurmV0041GetShares200Response) (map[string]*fairShareMetrics, error) {
	accounts := make(map[string]*fairShareMetrics)
	for _, s := range sharesResp.Shares.Shares {
		account := *s.Name
		if account == "root" {
			// we don't care about the root account
			continue
		}
		if _, exists := accounts[account]; !exists {
			accounts[account] = NewFairShareMetrics()
		}
		accounts[account].fairshare = GetFairShareEffectiveUsage(s)
	}
	return accounts, nil
}
