package api

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

func PopulateCache(ctx context.Context) error {
	var data []byte
	var err error

	cache := ctx.Value(types.ApiCacheKey).(*cache.Cache)

	numCalls := 5
	var wg sync.WaitGroup
	wg.Add(numCalls) // 5 different requests from slurm
	errors := make(chan error, numCalls)

	go func() {
		defer wg.Done()
		data, err = GetSlurmRestDiagResponse(ctx)
		if err != nil {
			errors <- fmt.Errorf("failed to get slurmrestd diagnostics response: %v", err)
		}
		cache.Set("diag", data, 0)
	}()

	go func() {
		defer wg.Done()
		data, err = GetSlurmRestNodesResponse(ctx)
		if err != nil {
			errors <- fmt.Errorf("failed to get slurmrestd nodes response: %v", err)
		}
		cache.Set("nodes", data, 0)
	}()

	go func() {
		defer wg.Done()
		data, err = GetSlurmRestJobsResponse(ctx)
		if err != nil {
			errors <- fmt.Errorf("failed to get slurmrestd jobs response: %v", err)
		}
		cache.Set("jobs", data, 0)
	}()

	go func() {
		defer wg.Done()
		data, err = GetSlurmRestPartitionsResponse(ctx)
		if err != nil {
			errors <- fmt.Errorf("failed to get slurmrestd partitions response: %v", err)
		}
		cache.Set("partitions", data, 0)
	}()

	go func() {
		defer wg.Done()
		data, err = GetSlurmRestSharesResponse(ctx)
		if err != nil {
			errors <- fmt.Errorf("failed to get slurmrestd shares response: %v", err)
		}
		cache.Set("shares", data, 0)
	}()

	slog.Info("waiting for workers")
	wg.Wait()

	for err := range errors {
		// yes i know it will only get the first error but it's almost certainly
		// going to be the same error 5 times
		return fmt.Errorf("errors encountered calling slurm api: %v", err)
	}

	return nil
}

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
