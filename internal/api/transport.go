package api

import (
	"context"
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

// GetSlurmRestResponse retrieves response data from slurm api
func GetSlurmRestResponse(ctx context.Context, endpointCtxKey types.Key) ([]byte, error) {
	var endpointStr string
	switch endpointCtxKey {
	case types.ApiDiagEndpointKey:
		endpointStr = "diag"
	case types.ApiJobsEndpointKey:
		endpointStr = "jobs"
	case types.ApiNodesEndpointKey:
		endpointStr = "nodes"
	case types.ApiPartitionsEndpointKey:
		endpointStr = "partitions"
	case types.ApiSharesEndpointKey:
		endpointStr = "shares"
	default:
		return nil, fmt.Errorf("invalid endpoint key")
	}
	nr, err := newSlurmRestRequest(ctx, endpointCtxKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new slurm rest request: %v", err)
	}
	resp, err := nr.Send()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve slurm rest response: %v", err)
	}
	if resp.StatusCode != 200 {
		slog.Debug("incorrect response status code", "endpoint", endpointStr, "code", resp.StatusCode, "body", string(resp.Body))
		return nil, fmt.Errorf("received incorrect status code for %s data", endpointStr)
	}
	return resp.Body, nil
}

// newSlurmRestRequest returns a new slurmRestRequest object which is used to perform
// http interactions with the slurmrest server. It configures everything up until
// the request is actually sent to get data.
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
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-SLURM-USER-NAME", apiUser)
	req.Header.Set("X-SLURM-USER-TOKEN", apiToken)

	return &slurmRestRequest{
		req:    req,
		client: &http.Client{},
	}, nil
}

// slurmRestRequest.Send is used to perform the request against the slurmrest
// server. It returns a *SlurmRestResponse which is a struct containing the
// response status code and the bytes of the response body.
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
