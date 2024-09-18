package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseSchedulerMetrics(t *testing.T) {
	diagBytes := util.ReadTestDataBytes("V0040OpenapiDiagResp.json")
	diagResp, _ := api.UnmarshalDiagResponse(diagBytes)
	_, err := ParseSchedulerMetrics(*diagResp)
	if err != nil {
		t.Fatalf("failed to parse scheduler metrics: %v\n", err)
	}
}
