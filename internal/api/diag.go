package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

// GetSlurmRestDiagResponse retrieves the diagnostic data respose from slurm api
func GetSlurmRestDiagResponse(ctx context.Context) ([]byte, error) {
	resp, err := newSlurmGETRequest(ctx, types.ApiDiagEndpointKey)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request for diag data: %v", err)
	}
	if resp.StatusCode != 200 {
		slog.Debug("incorrect status code for diag data", "code", resp.StatusCode, "body", string(resp.Body))
		return nil, fmt.Errorf("received incorrect status code for diag data")
	}
	return resp.Body, nil
}

// UnmarshalDiagResponse converts the response bytes into a slurm type
func UnmarshalDiagResponse(b []byte) (*types.V0040OpenapiDiagResp, error) {
	var diagResp types.V0040OpenapiDiagResp
	err := json.Unmarshal(b, &diagResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall diag response data: %v", err)
	}
	return &diagResp, nil
}
