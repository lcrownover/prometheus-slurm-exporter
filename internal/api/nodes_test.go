package api

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestUnmarshalNodesResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	_, err := UnmarshalNodesResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal nodes response: %v\n", err)
	}
}
