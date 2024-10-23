package api

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

// PopulateCache is used to populate the cache with data from the slurm api
func PopulateCache(ctx context.Context) error {
	slog.Debug("populating cache")
	var data []byte
	var err error

	apiCache := ctx.Value(types.ApiCacheKey).(*cache.Cache)

	var wg sync.WaitGroup
	wg.Add(len(endpoints))
	errors := make(chan error, len(endpoints))

	for _, e := range endpoints {
		go func(e endpoint) {
			defer wg.Done()
			data, err = GetSlurmRestResponse(ctx, e.key)
			if err != nil {
				errors <- fmt.Errorf("failed to get slurmrestd %s response: %v", e.path, err)
			}
			apiCache.Set(e.name, data, 0)
		}(e)
	}

	wg.Wait()
	close(errors)

	var errmsgs []string
	for err := range errors {
		errmsgs = append(errmsgs, err.Error())
		return fmt.Errorf("error(s) encountered calling slurm api: [%s]", strings.Join(errmsgs, ", "))
	}

	slog.Debug("finished populating cache")

	return nil
}

func WipeCache(ctx context.Context) error {
	apiCache := ctx.Value(types.ApiCacheKey).(*cache.Cache)
	apiCache.Delete("diag")
	apiCache.Delete("nodes")
	apiCache.Delete("jobs")
	apiCache.Delete("partitions")
	apiCache.Delete("shares")
	return nil
}
