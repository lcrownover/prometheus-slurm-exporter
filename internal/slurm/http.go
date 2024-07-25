package slurm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

// CleanseBaseURL removes any unneccessary elements from the provided url,
// such as "http://", etc
func CleanseBaseURL(url string) string {
	url = strings.ReplaceAll(url, "http://", "")
	url = strings.ReplaceAll(url, "https://", "")
	return url
}

type slurmRestRequest struct {
	req    *http.Request
	client *http.Client
}

type SlurmRestResponse struct {
	StatusCode int
	Body       []byte
}

func newSlurmRestRequest(ctx context.Context, k types.Key) (*slurmRestRequest, error) {
	apiUser := ctx.Value(types.ApiUserKey).(string)
	apiToken := ctx.Value(types.ApiTokenKey).(string)
	apiURL := ctx.Value(types.ApiURLKey).(string)
	apiEndpoint := ctx.Value(k).(string)

	url := fmt.Sprintf("http://%s/%s", apiURL, apiEndpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-SLURM-USER-NAME", apiUser)
	req.Header.Set("X-SLURM-USER-TOKEN", apiToken)

	return &slurmRestRequest{
		req:    req,
		client: &http.Client{},
	}, nil
}

func (sr slurmRestRequest) Send() (*SlurmRestResponse, error) {
	resp, err := sr.client.Do(sr.req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	sresp := SlurmRestResponse{}
	sresp.StatusCode = resp.StatusCode
	sresp.Body = body

	return &sresp, nil
}

// newGETRequest is a wrapper for net/http NewRequest so you only have to pass
// the endpoint. This packages the headers and client.
func newSlurmGETRequest(ctx context.Context, endpointCtxKey types.Key) (*SlurmRestResponse, error) {
	nr, err := newSlurmRestRequest(ctx, endpointCtxKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new slurm rest request: %v", err)
	}
	resp, err := nr.Send()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve slurm rest response: %v", err)
	}
	return resp, nil
}

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
