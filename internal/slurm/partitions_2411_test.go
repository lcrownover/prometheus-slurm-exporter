//go:build 2405

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestParsePartitionsMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesData, err := api.ProcessNodesResponse(nodesBytes)
	if err != nil {
		t.Fatalf("failed to extract nodes response: %v", err)
	}
	jobsBytes := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsData, _ := api.ProcessJobsResponse(jobsBytes)
	if err != nil {
		t.Fatalf("failed to extract jobs response: %v", err)
	}
	partitionsBytes := util.ReadTestDataBytes("V0041OpenapiPartitionResp.json")
	partitionData, _ := api.ProcessPartitionsResponse(partitionsBytes)
	if err != nil {
		t.Fatalf("failed to extract partitions response: %v", err)
	}
	data, err := ParsePartitionsMetrics(partitionData, jobsData, nodesData)
	if err != nil {
		t.Fatalf("failed to parse partitions metrics: %v", err)
	}
	tt := []struct {
		name    string
		metrics partitionMetrics
	}{
		{"name", partitionMetrics{0, 0, 1, 1, 0}},
		{"partition", partitionMetrics{0, 0, 0, 0, 2}},
		{"partitions", partitionMetrics{32, 36, -68, 0, 0}},
	}
	for _, tc := range tt {
		if data[tc.name].cpus_allocated != tc.metrics.cpus_allocated {
			t.Fatalf("expected %v, got %v", tc.metrics.cpus_allocated, data[tc.name].cpus_allocated)
		}
		if data[tc.name].cpus_idle != tc.metrics.cpus_idle {
			t.Fatalf("expected %v, got %v", tc.metrics.cpus_idle, data[tc.name].cpus_idle)
		}
		if data[tc.name].cpus_other != tc.metrics.cpus_other {
			t.Fatalf("expected %v, got %v", tc.metrics.cpus_other, data[tc.name].cpus_other)
		}
		if data[tc.name].cpus_total != tc.metrics.cpus_total {
			t.Fatalf("expected %v, got %v", tc.metrics.cpus_total, data[tc.name].cpus_total)
		}
		if data[tc.name].jobs_pending != tc.metrics.jobs_pending {
			t.Fatalf("expected %v, got %v", tc.metrics.jobs_pending, data[tc.name].jobs_pending)
		}
	}
}
