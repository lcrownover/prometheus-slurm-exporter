//go:build 2311

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestGetJobAccountName(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(fb)
	for _, j := range jobsResp.Jobs {
		_, err := GetJobAccountName(j)
		if err != nil {
			t.Fatalf("failed to get job account name: %v\n", err)
		}
	}
}

func TestGetJobState(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(fb)
	for _, j := range jobsResp.Jobs {
		_, err := GetJobState(j)
		if err != nil {
			t.Fatalf("failed to get job state: %v\n", err)
		}
	}
}

func TestGetJobCPUs(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(fb)
	for _, j := range jobsResp.Jobs {
		_, err := GetJobCPUs(j)
		if err != nil {
			t.Fatalf("failed to get job cpus: %v\n", err)
		}
	}
}

func TestGetNodeStates(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	for _, n := range nodesResp.Nodes {
		_, err := GetNodeStates(n)
		if err != nil {
			t.Fatalf("failed to get node states: %v\n", err)
		}
	}
}

func TestGetNodeGPUTotal(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	for _, n := range nodesResp.Nodes {
		_, err := GetNodeGPUTotal(n)
		if err != nil {
			t.Fatalf("failed to get node gpu total: %v\n", err)
		}
	}
}

func TestGetNodeGPUAllocated(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	for _, n := range nodesResp.Nodes {
		_, err := GetNodeGPUAllocated(n)
		if err != nil {
			t.Fatalf("failed to get node gpu allocated: %v\n", err)
		}
	}
}

func TestParseAccountMetrics(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(fb)
	_, err := ParseAccountsMetrics(jobsResp.Jobs)
	if err != nil {
		t.Fatalf("failed to parse account metrics: %v\n", err)
	}
}

func TestParseCPUsMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	nodesBytes := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(jobsBytes)
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseCPUsMetrics(*nodesResp, *jobsResp)
	if err != nil {
		t.Fatalf("failed to parse cpu metrics: %v\n", err)
	}
}

func TestParseGPUsMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseGPUsMetrics(*nodesResp)
	if err != nil {
		t.Fatalf("failed to parse gpu metrics: %v\n", err)
	}
}

func TestParseNodeMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseNodeMetrics(*nodesResp)
	if err != nil {
		t.Fatalf("failed to parse nodes metrics: %v\n", err)
	}
}

func TestParseNodesMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseNodesMetrics(*nodesResp)
	if err != nil {
		t.Fatalf("failed to parse nodes metrics: %v\n", err)
	}
}

func TestParsePartitionsMetrics(t *testing.T) {
	partitionsBytes := util.ReadTestDataBytes("V0040OpenapiPartitionResp.json")
	partitionResp, _ := api.UnmarshalPartitionsResponse(partitionsBytes)
	_, err := ParsePartitionsMetrics(*partitionResp)
	if err != nil {
		t.Fatalf("failed to parse partitions metrics: %v\n", err)
	}
}

func TestParseQueueMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(jobsBytes)
	_, err := ParseQueueMetrics(*jobsResp)
	if err != nil {
		t.Fatalf("failed to parse queue metrics: %v\n", err)
	}
}

func TestParseSchedulerMetrics(t *testing.T) {
	diagBytes := util.ReadTestDataBytes("V0040OpenapiDiagResp.json")
	diagResp, _ := api.UnmarshalDiagResponse(diagBytes)
	_, err := ParseSchedulerMetrics(*diagResp)
	if err != nil {
		t.Fatalf("failed to parse scheduler metrics: %v\n", err)
	}
}

func TestParseSharesMetrics(t *testing.T) {
	sharesBytes := util.ReadTestDataBytes("V0041OpenapiSharesResp.json")
	sharesResp, _ := api.UnmarshalSharesResponse(sharesBytes)
	_, err := ParseFairShareMetrics(*sharesResp)
	if err != nil {
		t.Fatalf("failed to parse fair share metrics: %v\n", err)
	}
}

func TestParseUsersMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(jobsBytes)
	_, err := ParseUsersMetrics(*jobsResp)
	if err != nil {
		t.Fatalf("failed to parse users metrics: %v\n", err)
	}
}
