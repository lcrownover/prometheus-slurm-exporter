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
		slog.Info("setting diag cache data", "data", data)
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
	close(errors)
	slog.Info("done waiting for workers")

	for err := range errors {
		// yes i know it will only get the first error but it's almost certainly
		// going to be the same error 5 times
		return fmt.Errorf("errors encountered calling slurm api: %v", err)
	}

	return nil
}
