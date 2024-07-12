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
	return rc.jobsCache
}

type jobsCache struct {
	ctx            context.Context
	data           *types.V0040OpenapiJobInfoResp
	expiration     int64
	lock           *sync.Mutex
	timeoutSeconds int
}

func newJobsCache(ctx context.Context, timeoutSeconds int) *jobsCache {
	return &jobsCache{ctx: ctx, data: nil, expiration: util.NowEpoch(), lock: &sync.Mutex{}, timeoutSeconds: timeoutSeconds}
}

func (jc *jobsCache) Jobs() []types.V0040JobInfo {
	jc.Refresh()
	return jc.data.Jobs
}

func (jc *jobsCache) Expiration() int64 {
	return jc.expiration
}

func (jc *jobsCache) Refresh() error {
	if !IsExpired(jc, jc.timeoutSeconds) {
		return nil
	}
	resp, err := util.NewSlurmGETRequest(jc.ctx, types.ApiJobsEndpointKey)
	if err != nil {
		return fmt.Errorf("failed to perform get request for job data: %v", err)
	}
	if resp.StatusCode != 200 {
		slog.Debug("incorrect status code for job data", "code", resp.StatusCode, "body", string(resp.Body))
		return fmt.Errorf("received incorrect status code for job data")
	}
	var jobsResp types.V0040OpenapiJobInfoResp
	err = json.Unmarshal(resp.Body, &jobsResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshall job response data: %v", err)
	}
	jc.data = &jobsResp
	jc.expiration = util.NowEpoch()
	return nil
}
