//go:build 2311

package api

import (
	"encoding/json"
	"fmt"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

// UnmarshalDiagResponse converts the response bytes into a slurm type
func UnmarshalDiagResponse(b []byte) (*types.V0040OpenapiDiagResp, error) {
	var diagResp types.V0040OpenapiDiagResp
	err := json.Unmarshal(b, &diagResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall diag response data: %v", err)
	}
	return &diagResp, nil
}

// UnmarshalJobsResponse converts the response bytes into a slurm type
func UnmarshalJobsResponse(b []byte) (*types.V0040OpenapiJobInfoResp, error) {
	var jobsResp types.V0040OpenapiJobInfoResp
	err := json.Unmarshal(b, &jobsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall job response data: %v", err)
	}
	return &jobsResp, nil
}

// UnmarshalNodesResponse converts the response bytes into a slurm type
func UnmarshalNodesResponse(b []byte) (*types.V0040OpenapiNodesResp, error) {
	var nodesResp types.V0040OpenapiNodesResp
	err := json.Unmarshal(b, &nodesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall nodes response data: %v", err)
	}
	return &nodesResp, nil
}

// UnmarshalPartitionsResponse converts the response bytes into a slurm type
func UnmarshalPartitionsResponse(b []byte) (*types.V0040OpenapiPartitionResp, error) {
	var partitionsResp types.V0040OpenapiPartitionResp
	err := json.Unmarshal(b, &partitionsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall partitions response data: %v", err)
	}
	return &partitionsResp, nil
}

// UnmarshalSharesResponse converts the response bytes into a slurm type
func UnmarshalSharesResponse(b []byte) (*types.V0040OpenapiSharesResp, error) {
	var sharesResp types.V0040OpenapiSharesResp
	err := json.Unmarshal(b, &sharesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall shares response data: %v", err)
	}
	return &sharesResp, nil
}
