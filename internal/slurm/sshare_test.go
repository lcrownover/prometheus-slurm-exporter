package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseSharesMetrics(t *testing.T) {
	sharesBytes := util.ReadTestDataBytes("V0040OpenapiSharesResp.json")
	sharesResp, _ := api.UnmarshalSharesResponse(sharesBytes)
	_, err := ParseFairShareMetrics(*sharesResp)
	if err != nil {
		t.Fatalf("failed to parse fair share metrics: %v\n", err)
	}
}
