package api

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestUnmarshalDiagResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiDiagResp.json")
	_, err := UnmarshalDiagResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal diag response: %v\n", err)
	}
}
