package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

// GetSlurmRestSharesResponse retrieves the fair share response from slurm api
func GetSlurmRestSharesResponse(ctx context.Context) ([]byte, error) {
	resp, err := newSlurmGETRequest(ctx, types.ApiSharesEndpointKey)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request for shares data: %v", err)
	}
	if resp.StatusCode != 200 {
		slog.Debug("incorrect status code for shares data", "code", resp.StatusCode, "body", string(resp.Body))
		return nil, fmt.Errorf("received incorrect status code for shares data")
	}
	return resp.Body, nil
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
