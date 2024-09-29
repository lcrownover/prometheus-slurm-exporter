//go:build 2405

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseCPUsMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(jobsBytes)
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	data, err := ParseCPUsMetrics(*nodesResp, *jobsResp)
	if err != nil {
		t.Fatalf("failed to parse cpu metrics: %v", err)
	}
	tt := []cpusMetrics{
		{0, 0, 18, 18},
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
	}
}
