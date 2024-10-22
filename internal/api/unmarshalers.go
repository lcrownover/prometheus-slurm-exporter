package api

import (
	"encoding/json"
	"fmt"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func ProcessDiagResponse(b []byte) (*DiagData, error) {
	var r DiagResp
	err := json.Unmarshal(b, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall diag response data: %v", err)
	}
	d := NewDiagData()
	d.FromResponse(r)
	return d, nil
}

// ProcessJobsResponse converts the response bytes into a slurm type
func ProcessJobsResponse(b []byte) (*JobsData, error) {
	var r JobsResp
	err := json.Unmarshal(b, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall job response data: %v", err)
	}
	d := NewJobsData()
	d.FromResponse(r)
	return d, nil
}

// ProcessNodesResponse converts the response bytes into a slurm type
func ProcessNodesResponse(b []byte) (*NodesData, error) {
	var r NodesResp
	err := json.Unmarshal(b, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall nodes response data: %v", err)
	}
	d := NewNodesData()
	d.FromResponse(r)
	return d, nil
}

// ProcessPartitionsResponse converts the response bytes into a slurm type
func ProcessPartitionsResponse(b []byte) (*PartitionsData, error) {
	var r PartitionsResp

	err := json.Unmarshal(b, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall partitions response data: %v", err)
	}
	d := NewPartitionsData()
	d.FromResponse(r)
	return d, nil
}

// ProcessSharesResponse converts the response bytes into a slurm type
func ProcessSharesResponse(b []byte) (*SharesData, error) {
	b = util.CleanseInfinity(b)
	var r SharesResp
	err := json.Unmarshal(b, &r)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall shares response data: %v", err)
	}

	d := NewSharesData()
	d.FromResponse(r)
	return d, nil
}