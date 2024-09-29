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
