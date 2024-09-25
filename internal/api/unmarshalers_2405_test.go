//go:build 2405

package api

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestUnmarshalDiagResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("SlurmV0041GetDiag200Response.json")
	_, err := UnmarshalDiagResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal diag response: %v\n", err)
	}
}

func TestUnmarshalJobsResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiJobInfoResp.json")
	_, err := UnmarshalJobsResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal jobs response: %v\n", err)
	}
}

func TestUnmarshalNodesResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiNodesResp.json")
	_, err := UnmarshalNodesResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal nodes response: %v\n", err)
	}
}

func TestUnmarshalPartitionsResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiPartitionResp.json")
	_, err := UnmarshalPartitionsResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal partition response: %v\n", err)
	}
}

func TestUnmarshalSharesResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0041OpenapiSharesResp.json")
	_, err := UnmarshalSharesResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal shares response: %v\n", err)
	}
}
