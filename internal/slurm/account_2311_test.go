//go:build 2311

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParseAccountsMetrics(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsData, _ := api.ExtractJobsData(fb)
	data, err := ParseAccountsMetrics(*jobsData)
	if err != nil {
		t.Fatalf("failed to parse accounts metrics: %v", err)
	}
	tt := []struct {
		account string
		metrics JobMetrics
	}{
		{"account", JobMetrics{2, 14, 0, 0, 0}},
	}
	for _, tc := range tt {
		if data[tc.account].pending != tc.metrics.pending {
			t.Fatalf("expected %v, got %v", tc.metrics.pending, data[tc.account].pending)
		}
		if data[tc.account].pending_cpus != tc.metrics.pending_cpus {
			t.Fatalf("expected %v, got %v", tc.metrics.pending_cpus, data[tc.account].pending_cpus)
		}
		if data[tc.account].running != tc.metrics.running {
			t.Fatalf("expected %v, got %v", tc.metrics.running, data[tc.account].running)
		}
		if data[tc.account].running_cpus != tc.metrics.running_cpus {
			t.Fatalf("expected %v, got %v", tc.metrics.running_cpus, data[tc.account].running_cpus)
		}
		if data[tc.account].suspended != tc.metrics.suspended {
			t.Fatalf("expected %v, got %v", tc.metrics.suspended, data[tc.account].suspended)
		}
	}
}
