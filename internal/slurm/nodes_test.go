package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseNodesMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseNodesMetrics(*nodesResp)
	if err != nil {
		t.Fatalf("failed to parse nodes metrics: %v\n", err)
	}
}
