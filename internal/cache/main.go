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
	lock.Lock()
	defer lock.Unlock()
	// cache exists, return it
	if singleResponseCache != nil {
		return singleResponseCache
	}
	// create a new cache if there's not a global one
	return newResponseCache(ctx)
}

// For each type of slurmrestd request, we add a pointer to the request cache
// This cache object is a singleton shared across all the slurm module code
// so we don't make multiple requests to the same endpoint for different metrics
type responseCache struct {
	ctx       context.Context
	jobsCache *jobsCache
}

func newResponseCache(ctx context.Context) *responseCache {
	return &responseCache{
		ctx:       ctx,
		jobsCache: newJobsCache(ctx),
	}
}

type RefreshableCache interface {
	LastRefresh() int64
	Ctx() context.Context
	Refresh() error
}

func TimeoutSeconds(ctx context.Context) int64 {
	var timeoutSeconds int64
	timeoutSecondsVal := ctx.Value(types.CacheTimeoutSecondsKey)
	if timeoutSecondsVal == nil {
		slog.Error("timeout seconds should not be nil. report to developer")
		timeoutSeconds = 5
	} else {
		timeoutSeconds = int64(timeoutSecondsVal.(int))
	}
	slog.Debug("timeout seconds", "value", timeoutSeconds)
	return int64(timeoutSeconds)
}

func IsExpired[T RefreshableCache](item T) bool {
	return time.Now().Unix() > item.LastRefresh()+TimeoutSeconds(item.Ctx())
}
