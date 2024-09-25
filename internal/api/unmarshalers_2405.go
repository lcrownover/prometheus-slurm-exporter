//go:build 2405

package api

import (
	"encoding/json"
	"fmt"

	openapi "github.com/lcrownover/openapi-slurm-24-05"
)

// UnmarshalDiagResponse converts the response bytes into a slurm type
func UnmarshalDiagResponse(b []byte) (*openapi.SlurmV0041GetDiag200Response, error) {
	var diagResp openapi.SlurmV0041GetDiag200Response
	err := json.Unmarshal(b, &diagResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall diag response data: %v", err)
	}
	return &diagResp, nil
}

// UnmarshalJobsResponse converts the response bytes into a slurm type
func UnmarshalJobsResponse(b []byte) (*openapi.V0041OpenapiJobInfoResp, error) {
	var jobsResp openapi.V0041OpenapiJobInfoResp
	err := json.Unmarshal(b, &jobsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall job response data: %v", err)
	}
	return &jobsResp, nil
}

// UnmarshalNodesResponse converts the response bytes into a slurm type
func UnmarshalNodesResponse(b []byte) (*openapi.V0041OpenapiNodesResp, error) {
	var nodesResp openapi.V0041OpenapiNodesResp
	err := json.Unmarshal(b, &nodesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall nodes response data: %v", err)
	}
	return &nodesResp, nil
}

// UnmarshalPartitionsResponse converts the response bytes into a slurm type
func UnmarshalPartitionsResponse(b []byte) (*openapi.V0041OpenapiPartitionResp, error) {
	var partitionsResp openapi.V0041OpenapiPartitionResp
	err := json.Unmarshal(b, &partitionsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall partitions response data: %v", err)
	}
	return &partitionsResp, nil
}

// UnmarshalSharesResponse converts the response bytes into a slurm type
func UnmarshalSharesResponse(b []byte) (*openapi.SlurmV0041GetShares200Response, error) {
	var sharesResp openapi.SlurmV0041GetShares200Response
	err := json.Unmarshal(b, &sharesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall shares response data: %v", err)
	}
	return &sharesResp, nil
}
