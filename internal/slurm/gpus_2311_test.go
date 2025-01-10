//go:build 2311

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseGPUsMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesData, _ := api.ProcessNodesResponse(nodesBytes)
	data, err := ParseGPUsMetrics(nodesData)
	if err != nil {
		t.Fatalf("failed to parse gpu metrics: %v", err)
	}
	tt := []gpusMetrics{
		{0, 0, 0, 0, 0},
	}
	for _, tc := range tt {
		if data.alloc != tc.alloc {
			t.Fatalf("expected %v, got %v", tc.alloc, data.alloc)
		}
		if data.idle != tc.idle {
			t.Fatalf("expected %v, got %v", tc.idle, data.idle)
		}
		if data.other != tc.other {
			t.Fatalf("expected %v, got %v", tc.other, data.other)
		}
		if data.total != tc.total {
			t.Fatalf("expected %v, got %v", tc.total, data.total)
		}
		if data.utilization != tc.utilization {
			t.Fatalf("expected %v, got %v", tc.utilization, data.utilization)
		}
	}
}
