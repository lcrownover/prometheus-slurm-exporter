package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseGPUsMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseGPUsMetrics(*nodesResp)
	if err != nil {
		t.Fatalf("failed to parse gpu metrics: %v\n", err)
	}
}
