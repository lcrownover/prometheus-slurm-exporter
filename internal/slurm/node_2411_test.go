//go:build 2411

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseNodeMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesData, _ := api.ProcessNodesResponse(nodesBytes)
	data, err := ParseNodeMetrics(nodesData)
	if err != nil {
		t.Fatalf("failed to parse nodes metrics: %v", err)
	}
	tt := []nodeMetrics{
		{6, 4, 8, 9, 0, 9, "invalid|invalid"},
	}
	for _, tc := range tt {
		if data["hostname"].memAlloc != tc.memAlloc {
			t.Fatalf("expected %v, got %v", tc.memAlloc, data["hostname"].memAlloc)
		}
		if data["hostname"].memTotal != tc.memTotal {
			t.Fatalf("expected %v, got %v", tc.memTotal, data["hostname"].memTotal)
		}
		if data["hostname"].cpuAlloc != tc.cpuAlloc {
			t.Fatalf("expected %v, got %v", tc.cpuAlloc, data["hostname"].cpuAlloc)
		}
		if data["hostname"].cpuIdle != tc.cpuIdle {
			t.Fatalf("expected %v, got %v", tc.cpuIdle, data["hostname"].cpuIdle)
		}
		if data["hostname"].cpuOther != tc.cpuOther {
			t.Fatalf("expected %v, got %v", tc.cpuOther, data["hostname"].cpuOther)
		}
		if data["hostname"].cpuTotal != tc.cpuTotal {
			t.Fatalf("expected %v, got %v", tc.cpuTotal, data["hostname"].cpuTotal)
		}
	}
}
