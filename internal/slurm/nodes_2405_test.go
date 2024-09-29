//go:build 2405

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseNodesMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	data, err := ParseNodesMetrics(*nodesResp)
	if err != nil {
		t.Fatalf("failed to parse nodes metrics: %v", err)
	}
	tt := []nodesMetrics{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	for _, tc := range tt {
		if data.alloc != tc.alloc {
			t.Fatalf("expected %v, got %v", tc.alloc, data.alloc)
		}
		if data.comp != tc.comp {
			t.Fatalf("expected %v, got %v", tc.comp, data.comp)
		}
		if data.down != tc.down {
			t.Fatalf("expected %v, got %v", tc.down, data.down)
		}
		if data.drain != tc.drain {
			t.Fatalf("expected %v, got %v", tc.drain, data.drain)
		}
		if data.err != tc.err {
			t.Fatalf("expected %v, got %v", tc.err, data.err)
		}
		if data.fail != tc.fail {
			t.Fatalf("expected %v, got %v", tc.fail, data.fail)
		}
		if data.idle != tc.idle {
			t.Fatalf("expected %v, got %v", tc.idle, data.idle)
		}
		if data.maint != tc.maint {
			t.Fatalf("expected %v, got %v", tc.maint, data.maint)
		}
		if data.mix != tc.mix {
			t.Fatalf("expected %v, got %v", tc.mix, data.mix)
		}
		if data.resv != tc.resv {
			t.Fatalf("expected %v, got %v", tc.resv, data.resv)
		}
	}
}
