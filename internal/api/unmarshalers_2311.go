//go:build 2311

package api

import (
	"encoding/json"
	"fmt"
	"strings"

	openapi "github.com/lcrownover/openapi-slurm-23-11"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

var apiVersion = "23.11"

// UnmarshalDiagResponse converts the response bytes into a slurm type
func UnmarshalDiagResponse(b []byte) (*openapi.V0040OpenapiDiagResp, error) {
	var diagResp openapi.V0040OpenapiDiagResp
	err := json.Unmarshal(b, &diagResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall diag response data: %v", err)
	}
	return &diagResp, nil
}

// UnmarshalJobsResponse converts the response bytes into a slurm type
func UnmarshalJobsResponse(b []byte) (*openapi.V0040OpenapiJobInfoResp, error) {
	var jobsResp openapi.V0040OpenapiJobInfoResp
	err := json.Unmarshal(b, &jobsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall job response data: %v", err)
	}
	return &jobsResp, nil
}

// UnmarshalNodesResponse converts the response bytes into a slurm type
func UnmarshalNodesResponse(b []byte) (*openapi.V0040OpenapiNodesResp, error) {
	var nodesResp openapi.V0040OpenapiNodesResp
	err := json.Unmarshal(b, &nodesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall nodes response data: %v", err)
	}
	return &nodesResp, nil
}

// UnmarshalPartitionsResponse converts the response bytes into a slurm type
func UnmarshalPartitionsResponse(b []byte) (*openapi.V0040OpenapiPartitionResp, error) {
	var partitionsResp openapi.V0040OpenapiPartitionResp
	err := json.Unmarshal(b, &partitionsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall partitions response data: %v", err)
	}
	return &partitionsResp, nil
}

// UnmarshalSharesResponse converts the response bytes into a slurm type
func UnmarshalSharesResponse(b []byte) (*openapi.V0040OpenapiSharesResp, error) {
	// this is disgusting but the response has values of "Infinity" which are
	// not json unmarshal-able, so I manually replace all the "Infinity"s with the correct
	// float64 value that represents Infinity.
	// this will be fixed in v0.0.42
	// https://support.schedmd.com/show_bug.cgi?id=20817
	//
	// https://github.com/lcrownover/prometheus-slurm-exporter/issues/8
	// also reported that folks are getting "inf" back, so I'll protect for that too
	sharesRespString := util.RemoveWhitespace(string(b))
	maxFloatStr := ":1.7976931348623157e+308"
	sharesRespString = strings.ReplaceAll(sharesRespString, ":Infinity", maxFloatStr) // replacing the longer strings first should prevent any partial replacements
	sharesRespString = strings.ReplaceAll(sharesRespString, ":infinity", maxFloatStr)
	sharesRespString = strings.ReplaceAll(sharesRespString, ":Inf", maxFloatStr) // sometimes it'd return "inf", so let's cover for that too.
	sharesRespString = strings.ReplaceAll(sharesRespString, ":inf", maxFloatStr)
	sharesRespBytes := []byte(sharesRespString)
	// end hack

	var sharesResp openapi.V0040OpenapiSharesResp
	err := json.Unmarshal(sharesRespBytes, &sharesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall shares response data: %v", err)
	}

	return &sharesResp, nil
}
