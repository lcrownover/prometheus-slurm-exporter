package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

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
	slog.Debug("performing rest request", "endpoint", endpointStr)
	nr, err := newSlurmRestRequest(ctx, endpointCtxKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new slurm rest request: %v", err)
	}
	resp, err := nr.Send()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve slurm rest response: %v", err)
	}
	// sometimes slurm fails to get stuff. we want to error here
	if resp.StatusCode == 500 {
		slog.Debug("incorrect response status code", "endpoint", endpointStr, "code", resp.StatusCode, "body", string(resp.Body))

		// try to unmarshal the api error and give a better log
		var aed APIErrorData
		var errStr string
		err := json.Unmarshal(resp.Body, &aed)
		if err != nil {
			errStr = "tried to get more data about the error but failed. try debug mode for more information"
		}
		errStr = aed.ToString()
		return nil, fmt.Errorf("internal server error (500) from slurm controller getting %s data: %s", endpointStr, errStr)
	}
	// unauthorized responses should say that
	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("unauthorized: invalid credentials")
	}
	// otherwise, it should be status 200, so this catches unsupported status codes
	if resp.StatusCode != 200 {
		slog.Debug("incorrect response status code", "endpoint", endpointStr, "code", resp.StatusCode, "body", string(resp.Body))
		return nil, fmt.Errorf("received incorrect status code for %s data", endpointStr)
	}
	slog.Debug("successfully queried slurm rest data", "endpoint", endpointStr)
	return resp.Body, nil
}

// newSlurmRestRequest returns a new slurmRestRequest object which is used to perform
// http interactions with the slurmrest server. It configures everything up until
// the request is actually sent to get data.
func newSlurmRestRequest(ctx context.Context, k types.Key) (*slurmRestRequest, error) {
	apiURL := ctx.Value(types.ApiURLKey).(string)

	if strings.HasPrefix(apiURL, "unix://") {
		return newSlurmUnixRestRequest(ctx, k)
	} else if strings.HasPrefix(apiURL, "http://") || strings.HasPrefix(apiURL, "https://") {
		return newSlurmInetRestRequest(ctx, k)
	}
	return nil, fmt.Errorf("invalid SLURM_EXPORTER_API_URL: %s", apiURL)
}

func newSlurmInetRestRequest(ctx context.Context, k types.Key) (*slurmRestRequest, error) {
	apiUser := ctx.Value(types.ApiUserKey).(string)
	apiToken := ctx.Value(types.ApiTokenKey).(string)
	apiURL := ctx.Value(types.ApiURLKey).(string)
	apiEndpoint := ctx.Value(k).(string)

	url := fmt.Sprintf("%s/%s", apiURL, apiEndpoint)
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

func newSlurmUnixRestRequest(ctx context.Context, k types.Key) (*slurmRestRequest, error) {
	apiURL := ctx.Value(types.ApiURLKey).(string)
	apiEndpoint := ctx.Value(k).(string)

	socketPath := strings.TrimPrefix(apiURL, "unix:")
	url := fmt.Sprintf("http://unix/%s", apiEndpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	return &slurmRestRequest{
		req: req,
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
		},
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
