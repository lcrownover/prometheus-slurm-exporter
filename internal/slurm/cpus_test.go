package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseCPUsMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	nodesBytes := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(jobsBytes)
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseCPUsMetrics(*nodesResp, *jobsResp)
	if err != nil {
		t.Fatalf("failed to parse cpu metrics: %v\n", err)
	}
}
