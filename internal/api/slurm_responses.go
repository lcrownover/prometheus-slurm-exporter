package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

// GetSlurmRestDiagResponse retrieves the diagnostic data respose from slurm api
func GetSlurmRestDiagResponse(ctx context.Context) ([]byte, error) {
	cache := ctx.Value(types.ApiCacheKey).(*cache.Cache)
	ct := ctx.Value(types.ApiCacheTimeoutKey).(time.Duration)
	data, found := cache.Get("diag")
	if !found {
		resp, err := newSlurmGETRequest(ctx, types.ApiDiagEndpointKey)
		if err != nil {
			return nil, fmt.Errorf("failed to perform get request for diag data: %v", err)
		}
		if resp.StatusCode != 200 {
			slog.Debug("incorrect status code for diag data", "code", resp.StatusCode, "body", string(resp.Body))
			return nil, fmt.Errorf("received incorrect status code for diag data")
		}
		cache.Set("diag", resp.Body, ct)
		data = resp.Body
	}
	return data.([]byte), nil
}

// GetSlurmRestJobsResponse retrieves response bytes from slurm REST api
func GetSlurmRestJobsResponse(ctx context.Context) ([]byte, error) {
	cache := ctx.Value(types.ApiCacheKey).(*cache.Cache)
	ct := ctx.Value(types.ApiCacheTimeoutKey).(time.Duration)
	data, found := cache.Get("jobs")
	if !found {
		resp, err := newSlurmGETRequest(ctx, types.ApiJobsEndpointKey)
		if err != nil {
			return nil, fmt.Errorf("failed to perform get request for job data: %v", err)
		}
		if resp.StatusCode != 200 {
			slog.Debug("incorrect status code for job data", "code", resp.StatusCode, "body", string(resp.Body))
			return nil, fmt.Errorf("received incorrect status code for job data")
		}
		cache.Set("jobs", resp.Body, ct)
		data = resp.Body
	}
	return data.([]byte), nil
}

// GetSlurmRestNodesResponse retrieves the list of nodes registered to slurm
func GetSlurmRestNodesResponse(ctx context.Context) ([]byte, error) {
	cache := ctx.Value(types.ApiCacheKey).(*cache.Cache)
	ct := ctx.Value(types.ApiCacheTimeoutKey).(time.Duration)
	data, found := cache.Get("nodes")
	if !found {
		resp, err := newSlurmGETRequest(ctx, types.ApiNodesEndpointKey)
		if err != nil {
			return nil, fmt.Errorf("failed to perform get request for node data: %v", err)
		}
		if resp.StatusCode != 200 {
			slog.Debug("incorrect status code for node data", "code", resp.StatusCode, "body", string(resp.Body))
			return nil, fmt.Errorf("received incorrect status code for node data")
		}
		cache.Set("nodes", resp.Body, ct)
		data = resp.Body
	} else {
		slog.Info("using cached nodes data")
	}
	return data.([]byte), nil
}

// GetSlurmRestPartitionsResponse retrieves response bytes from slurm REST api
func GetSlurmRestPartitionsResponse(ctx context.Context) ([]byte, error) {
	cache := ctx.Value(types.ApiCacheKey).(*cache.Cache)
	ct := ctx.Value(types.ApiCacheTimeoutKey).(time.Duration)
	data, found := cache.Get("partitions")
	if !found {
		resp, err := newSlurmGETRequest(ctx, types.ApiPartitionsEndpointKey)
		if err != nil {
			return nil, fmt.Errorf("failed to perform get request for partition data: %v", err)
		}
		if resp.StatusCode != 200 {
			slog.Debug("incorrect status code for partition data", "code", resp.StatusCode, "body", string(resp.Body))
			return nil, fmt.Errorf("received incorrect status code for partition data")
		}
		cache.Set("partitions", resp.Body, ct)
		data = resp.Body
	}
	return data.([]byte), nil
}

// GetSlurmRestSharesResponse retrieves the fair share response from slurm api
func GetSlurmRestSharesResponse(ctx context.Context) ([]byte, error) {
	cache := ctx.Value(types.ApiCacheKey).(*cache.Cache)
	ct := ctx.Value(types.ApiCacheTimeoutKey).(time.Duration)
	data, found := cache.Get("shares")
	if !found {
		resp, err := newSlurmGETRequest(ctx, types.ApiSharesEndpointKey)
		if err != nil {
			return nil, fmt.Errorf("failed to perform get request for shares data: %v", err)
		}
		if resp.StatusCode != 200 {
			slog.Debug("incorrect status code for shares data", "code", resp.StatusCode, "body", string(resp.Body))
			return nil, fmt.Errorf("received incorrect status code for shares data")
		}
		cache.Set("shares", resp.Body, ct)
		data = resp.Body
	}
	return data.([]byte), nil
}
