//go:build 2405

package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestGetJobAccountName(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(fb)
	for _, j := range jobsResp.Jobs {
		_, err := GetJobAccountName(j)
		if err != nil {
			t.Fatalf("failed to get job account name: %v", err)
		}
	}
}

func TestGetJobState(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(fb)
	for _, j := range jobsResp.Jobs {
		_, err := GetJobState(j)
		if err != nil {
			t.Fatalf("failed to get job state: %v", err)
		}
	}
}

func TestGetJobCPUs(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(fb)
	for _, j := range jobsResp.Jobs {
		_, err := GetJobCPUs(j)
		if err != nil {
			t.Fatalf("failed to get job cpus: %v", err)
		}
	}
}

func TestGetNodeStates(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesResp, err := api.UnmarshalNodesResponse(nodesBytes)
	if err != nil {
		t.Fatalf("failed to unmarshal nodes response: %v", err)
	}
	for _, n := range nodesResp.Nodes {
		_, err := GetNodeStates(n)
		if err != nil {
			t.Fatalf("failed to get node states: %v", err)
		}
	}
}

func TestGetNodeGPUTotal(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	for _, n := range nodesResp.Nodes {
		_, err := GetNodeGPUTotal(n)
		if err != nil {
			t.Fatalf("failed to get node gpu total: %v", err)
		}
	}
}

func TestGetNodeGPUAllocated(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	for _, n := range nodesResp.Nodes {
		_, err := GetNodeGPUAllocated(n)
		if err != nil {
			t.Fatalf("failed to get node gpu allocated: %v", err)
		}
	}
}

func TestParseAccountMetrics(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(fb)
	_, err := ParseAccountsMetrics(*jobsResp)
	if err != nil {
		t.Fatalf("failed to parse account metrics: %v", err)
	}
}

func TestParseCPUsMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(jobsBytes)
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseCPUsMetrics(*nodesResp, *jobsResp)
	if err != nil {
		t.Fatalf("failed to parse cpu metrics: %v", err)
	}
}

func TestParseGPUsMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseGPUsMetrics(*nodesResp)
	if err != nil {
		t.Fatalf("failed to parse gpu metrics: %v", err)
	}
}

func TestParseNodeMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseNodeMetrics(*nodesResp)
	if err != nil {
		t.Fatalf("failed to parse nodes metrics: %v", err)
	}
}

func TestParseNodesMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesResp, _ := api.UnmarshalNodesResponse(nodesBytes)
	_, err := ParseNodesMetrics(*nodesResp)
	if err != nil {
		t.Fatalf("failed to parse nodes metrics: %v", err)
	}
}

func TestParsePartitionsMetrics(t *testing.T) {
	nodesBytes := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	nodesResp, err := api.UnmarshalNodesResponse(nodesBytes)
	if err != nil {
		t.Fatalf("failed to unmarshal nodes response: %v", err)
	}
	jobsBytes := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(jobsBytes)
	if err != nil {
		t.Fatalf("failed to unmarshal jobs response: %v", err)
	}
	partitionsBytes := util.ReadTestDataBytes("V0041OpenapiPartitionResp.json")
	partitionResp, _ := api.UnmarshalPartitionsResponse(partitionsBytes)
	if err != nil {
		t.Fatalf("failed to unmarshal partitions response: %v", err)
	}
	_, err = ParsePartitionsMetrics(*partitionResp, *jobsResp, *nodesResp)
	if err != nil {
		t.Fatalf("failed to parse partitions metrics: %v", err)
	}
}

func TestParseQueueMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(jobsBytes)
	_, err := ParseQueueMetrics(*jobsResp)
	if err != nil {
		t.Fatalf("failed to parse queue metrics: %v", err)
	}
}

func TestParseSchedulerMetrics(t *testing.T) {
	diagBytes := util.ReadTestDataBytes("SlurmV0041GetDiag200Response.json")
	diagResp, _ := api.UnmarshalDiagResponse(diagBytes)
	_, err := ParseSchedulerMetrics(*diagResp)
	if err != nil {
		t.Fatalf("failed to parse scheduler metrics: %v", err)
	}
}

func TestParseSharesMetrics(t *testing.T) {
	sharesBytes := util.ReadTestDataBytes("SlurmV0041GetShares200Response.json")
	sharesResp, _ := api.UnmarshalSharesResponse(sharesBytes)
	_, err := ParseFairShareMetrics(*sharesResp)
	if err != nil {
		t.Fatalf("failed to parse fair share metrics: %v", err)
	}
}

func TestParseUsersMetrics(t *testing.T) {
	jobsBytes := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	jobsResp, _ := api.UnmarshalJobsResponse(jobsBytes)
	_, err := ParseUsersMetrics(*jobsResp)
	if err != nil {
		t.Fatalf("failed to parse users metrics: %v", err)
	}
}
