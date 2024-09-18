package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

// GetSlurmRestPartitionsResponse retrieves response bytes from slurm REST api
func GetSlurmRestPartitionsResponse(ctx context.Context) ([]byte, error) {
	resp, err := newSlurmGETRequest(ctx, types.ApiPartitionsEndpointKey)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request for partition data: %v", err)
	}
	if resp.StatusCode != 200 {
		slog.Debug("incorrect status code for partition data", "code", resp.StatusCode, "body", string(resp.Body))
		return nil, fmt.Errorf("received incorrect status code for partition data")
	}
	return resp.Body, nil
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
