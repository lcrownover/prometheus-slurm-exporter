//go:build 2405

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseAccountMetrics(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(fb)
	_, err := ParseAccountsMetrics(jobsResp.Jobs)
	if err != nil {
		t.Fatalf("failed to parse account metrics: %v\n", err)
	}
}
