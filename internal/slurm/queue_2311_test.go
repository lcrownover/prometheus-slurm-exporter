//go:build 2311

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseQueueMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsData, _ := api.ExtractJobsData(jobsBytes)
	data, err := ParseQueueMetrics(jobsData)
	if err != nil {
		t.Fatalf("failed to parse queue metrics: %v", err)
	}
	tt := []queueMetrics{
		{0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	for _, tc := range tt {
		if data.pending != tc.pending {
			t.Fatalf("expected %v, got %v", tc.pending, data.pending)
		}
		if data.pending_dep != tc.pending_dep {
			t.Fatalf("expected %v, got %v", tc.pending_dep, data.pending_dep)
		}
		if data.running != tc.running {
			t.Fatalf("expected %v, got %v", tc.running, data.running)
		}
		if data.suspended != tc.suspended {
			t.Fatalf("expected %v, got %v", tc.suspended, data.suspended)
		}
		if data.cancelled != tc.cancelled {
			t.Fatalf("expected %v, got %v", tc.cancelled, data.cancelled)
		}
		if data.completing != tc.completing {
			t.Fatalf("expected %v, got %v", tc.completing, data.completing)
		}
		if data.completed != tc.completed {
			t.Fatalf("expected %v, got %v", tc.completed, data.completed)
		}
		if data.configuring != tc.configuring {
			t.Fatalf("expected %v, got %v", tc.configuring, data.configuring)
		}
		if data.failed != tc.failed {
			t.Fatalf("expected %v, got %v", tc.failed, data.failed)
		}
		if data.timeout != tc.timeout {
			t.Fatalf("expected %v, got %v", tc.timeout, data.timeout)
		}
		if data.preempted != tc.preempted {
			t.Fatalf("expected %v, got %v", tc.preempted, data.preempted)
		}
		if data.node_fail != tc.node_fail {
			t.Fatalf("expected %v, got %v", tc.node_fail, data.node_fail)
		}
	}
}
