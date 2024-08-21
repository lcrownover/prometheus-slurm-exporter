package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

// GetSlurmRestNodesResponse retrieves the list of nodes registered to slurm
func GetSlurmRestNodesResponse(ctx context.Context) ([]byte, error) {
	resp, err := newSlurmGETRequest(ctx, types.ApiNodesEndpointKey)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request for node data: %v", err)
	}
	if resp.StatusCode != 200 {
		slog.Debug("incorrect status code for node data", "code", resp.StatusCode, "body", string(resp.Body))
		return nil, fmt.Errorf("received incorrect status code for node data")
	}
	return resp.Body, nil
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
