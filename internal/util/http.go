package util

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

func newSlurmRestRequest(ctx context.Context, k types.Key) (*slurmRestRequest, error) {
	slog.Debug("newSlurmRestRequest k", "value", k)
	apiUser := ctx.Value(types.ApiUserKey).(string)
	apiToken := ctx.Value(types.ApiTokenKey).(string)
	apiURL := ctx.Value(types.ApiURLKey).(string)
	slog.Debug("value of that k", "value", ctx.Value(k))
	slog.Debug("ctx", "value", ctx)
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

// NewGETRequest is a wrapper for net/http NewRequest so you only have to pass
// the endpoint. This packages the headers and client.
func NewSlurmGETRequest(ctx context.Context, endpointCtxKey types.Key) (*SlurmRestResponse, error) {
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
