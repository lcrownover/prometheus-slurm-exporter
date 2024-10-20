//go:build 2405

package slurm

import (
	"fmt"
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseSchedulerMetrics(t *testing.T) {
	diagBytes := util.ReadTestDataBytes("SlurmV0041GetDiag200Response.json")
	diagData, _ := api.ExtractDiagData(diagBytes)
	data, err := ParseSchedulerMetrics(diagData)
	if err != nil {
		t.Fatalf("failed to parse scheduler metrics: %v", err)
	}
	tt := []schedulerMetrics{
		{5, 5, 9, 4, 1, 6, 0, 7, 0, 6, 5, 3},
	}
	for _, tc := range tt {
		if data.threads != tc.threads {
			t.Fatalf("expected threads %v, got %v", tc.threads, data.threads)
		}
		if data.queue_size != tc.queue_size {
			t.Fatalf("expected queue_size %v, got %v", tc.queue_size, data.queue_size)
		}
		if data.dbd_queue_size != tc.dbd_queue_size {
			t.Fatalf("expected dbd_queue_size %v, got %v", tc.dbd_queue_size, data.dbd_queue_size)
		}
		if data.last_cycle != tc.last_cycle {
			t.Fatalf("expected last_cycle %v, got %v", tc.last_cycle, data.last_cycle)
		}
		if data.mean_cycle != tc.mean_cycle {
			t.Fatalf("expected mean_cycle %v, got %v", tc.mean_cycle, data.mean_cycle)
		}
		if data.cycle_per_minute != tc.cycle_per_minute {
			t.Fatalf("expected cycle_per_minute %v, got %v", tc.cycle_per_minute, data.cycle_per_minute)
		}
		if data.backfill_last_cycle != tc.backfill_last_cycle {
			t.Fatalf("expected backfill_last_cycle %v, got %v", tc.backfill_last_cycle, data.backfill_last_cycle)
		}
		if data.backfill_mean_cycle != tc.backfill_mean_cycle {
			t.Fatalf("expected backfill_mean_cycle %v, got %v", tc.backfill_mean_cycle, data.backfill_mean_cycle)
		}
		if data.backfill_depth_mean != tc.backfill_depth_mean {
			t.Fatalf("expected backfill_depth_mean %v, got %v", tc.backfill_depth_mean, data.backfill_depth_mean)
		}
		if data.total_backfilled_jobs_since_start != tc.total_backfilled_jobs_since_start {
			t.Fatalf("expected total_backfilled_jobs_since_start %v, got %v", tc.total_backfilled_jobs_since_start, data.total_backfilled_jobs_since_start)
		}
		if data.total_backfilled_jobs_since_cycle != tc.total_backfilled_jobs_since_cycle {
			t.Fatalf("expected total_backfilled_jobs_since_cycle %v, got %v", tc.total_backfilled_jobs_since_cycle, data.total_backfilled_jobs_since_cycle)
		}
		if data.total_backfilled_heterogeneous != tc.total_backfilled_heterogeneous {
			t.Fatalf("expected total_backfilled_heterogeneous %v, got %v", tc.total_backfilled_heterogeneous, data.total_backfilled_heterogeneous)
		}
	}
}
