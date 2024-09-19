package api

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestUnmarshalJobsResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	_, err := UnmarshalJobsResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal jobs response: %v\n", err)
	}
}

