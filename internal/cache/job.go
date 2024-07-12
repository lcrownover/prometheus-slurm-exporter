package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

// Accessor for the job cache
func (rc responseCache) JobsCache() *jobsCache {
	rc.jobs.Get()
	return rc.jobs
}

type jobsCache struct {
	ctx            context.Context
	Data           *types.V0040OpenapiJobInfoResp
	expiration     int64
	lock           *sync.Mutex
	timeoutSeconds int
}

func newJobsCache(ctx context.Context, timeoutSeconds int) *jobsCache {
	return &jobsCache{ctx: ctx, Data: nil, expiration: util.NowEpoch(), lock: &sync.Mutex{}, timeoutSeconds: timeoutSeconds}
}

func (ji jobsCache) Expiration() int64 {
	return ji.expiration
}

func (ji jobsCache) Get() (*types.V0040OpenapiJobInfoResp, error) {
	// return cached data if it's still good
	if !IsExpired(ji, ji.timeoutSeconds) {
		return ji.Data, nil
	}
	resp, err := util.NewSlurmGETRequest(ji.ctx, types.ApiJobsEndpointKey)
	if err != nil {
		return nil, fmt.Errorf("failed to perform get request for job data: %v", err)
	}
	if resp.StatusCode != 200 {
		slog.Debug("incorrect status code for job data", "code", resp.StatusCode, "body", string(resp.Body))
		return nil, fmt.Errorf("received incorrect status code for job data")
	}
	var sj types.V0040OpenapiJobInfoResp
	err = json.Unmarshal(resp.Body, &sj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall job response data: %v", err)
	}
	ji.Data = &sj
	ji.expiration = util.NowEpoch()
	return ji.Data, nil
}
