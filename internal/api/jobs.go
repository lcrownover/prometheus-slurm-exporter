package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

// GetSlurmRestJobsResponse retrieves response bytes from slurm REST api
func GetSlurmRestJobsResponse(ctx context.Context) ([]byte, error) {
	resp, err := newSlurmGETRequest(ctx, types.ApiJobsEndpointKey)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request for job data: %v", err)
	}
	if resp.StatusCode != 200 {
		slog.Debug("incorrect status code for job data", "code", resp.StatusCode, "body", string(resp.Body))
		return nil, fmt.Errorf("received incorrect status code for job data")
	}
	return resp.Body, nil
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
