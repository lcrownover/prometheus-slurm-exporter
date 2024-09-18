package api

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestUnmarshalSharesResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiSharesResp.json")
	_, err := UnmarshalSharesResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal shares response: %v\n", err)
	}
}
