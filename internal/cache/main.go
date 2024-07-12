package cache

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

var singleResponseCache *responseCache

var lock = &sync.Mutex{}

// GetResponseCache is a singleton generator for a requestCache
// Many of the slurm module functions need to share requests, so we cache the
// responses and share in a singleton
func GetResponseCache(ctx context.Context) *responseCache {
	var timeoutSeconds int
	timeoutSecondsVal := ctx.Value(types.CacheTimeoutSeconds)
	if timeoutSecondsVal == nil {
		slog.Error("timeoutSeconds should have been set in ctx. report to developer.")
		timeoutSeconds = 5
	}
	timeoutSeconds = timeoutSecondsVal.(int)
	lock.Lock()
	defer lock.Unlock()
	// cache exists, return it
	if singleResponseCache != nil {
		return singleResponseCache
	}
	// create a new cache if there's not a global one
	return newResponseCache(ctx, timeoutSeconds)
}

// For each type of slurmrestd request, we add a pointer to the request cache
// This cache object is a singleton shared across all the slurm module code
// so we don't make multiple requests to the same endpoint for different metrics
type responseCache struct {
	ctx            context.Context
	timeoutSeconds int
	jobsCache      *jobsCache
}

func newResponseCache(ctx context.Context, timeoutSeconds int) *responseCache {
	return &responseCache{
		ctx:            ctx,
		timeoutSeconds: timeoutSeconds,
		jobsCache:      newJobsCache(ctx, timeoutSeconds),
	}
}

type Expirable interface {
	Expiration() int64
}

func IsExpired[T Expirable](item T, timeoutSeconds int) bool {
	return time.Now().Unix() > item.Expiration()+int64(timeoutSeconds)
}
