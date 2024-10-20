//go:build 2405

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseUsersMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsData, _ := api.ExtractJobsData(jobsBytes)
	data, err := ParseUsersMetrics(jobsData)
	if err != nil {
		t.Fatalf("failed to parse users metrics: %v", err)
	}
	tt := []struct {
		userName string
		metrics  userJobMetrics
	}{
		{"user_name", userJobMetrics{2, 12, 0, 0, 0}},
	}
	for _, tc := range tt {
		if data[tc.userName].pending != tc.metrics.pending {
			t.Fatalf("expected pending %v, got %v", tc.metrics.pending, data[tc.userName].pending)
		}
		if data[tc.userName].pending_cpus != tc.metrics.pending_cpus {
			t.Fatalf("expected pending_cpus %v, got %v", tc.metrics.pending_cpus, data[tc.userName].pending_cpus)
		}
		if data[tc.userName].running != tc.metrics.running {
			t.Fatalf("expected running %v, got %v", tc.metrics.running, data[tc.userName].running)
		}
		if data[tc.userName].running_cpus != tc.metrics.running_cpus {
			t.Fatalf("expected running_cpus %v, got %v", tc.metrics.running_cpus, data[tc.userName].running_cpus)
		}
		if data[tc.userName].suspended != tc.metrics.suspended {
			t.Fatalf("expected suspended %v, got %v", tc.metrics.suspended, data[tc.userName].suspended)
		}
	}
}
