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
