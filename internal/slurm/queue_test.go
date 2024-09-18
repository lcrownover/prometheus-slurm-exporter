package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseQueueMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(jobsBytes)
	_, err := ParseQueueMetrics(*jobsResp)
	if err != nil {
		t.Fatalf("failed to parse queue metrics: %v\n", err)
	}
}
