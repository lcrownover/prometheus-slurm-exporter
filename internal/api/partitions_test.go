package api

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestUnmarshalPartitionsResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiPartitionsResp.json")
	_, err := UnmarshalPartitionsResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal partition response: %v\n", err)
	}
}
