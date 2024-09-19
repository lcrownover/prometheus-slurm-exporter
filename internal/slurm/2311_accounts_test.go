//go:build 2311

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseAccountMetrics(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(fb)
	_, err := ParseAccountsMetrics(jobsResp.Jobs)
	if err != nil {
		t.Fatalf("failed to parse account metrics: %v\n", err)
	}
}
