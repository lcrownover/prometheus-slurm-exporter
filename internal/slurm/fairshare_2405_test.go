//go:build 2405

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseSharesMetrics(t *testing.T) {
	sharesBytes := util.ReadTestDataBytes("SlurmV0041GetShares200Response.json")
	sharesResp, _ := api.UnmarshalSharesResponse(sharesBytes)
	data, err := ParseFairShareMetrics(*sharesResp)
	if err != nil {
		t.Fatalf("failed to parse fair share metrics: %v", err)
	}
	tt := []fairShareMetrics{
		{2.3021358869347655},
	}
	for _, tc := range tt {
		if data["name"].fairshare != tc.fairshare {
			t.Fatalf("expected %v, got %v", tc.fairshare, data["name"].fairshare)
		}
	}
}
