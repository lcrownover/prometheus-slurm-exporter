//go:build 2311

package slurm

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/akyoto/cache"
	openapi "github.com/lcrownover/openapi-slurm-23-11"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

type UsersCollector struct {
	ctx          context.Context
	pending      *prometheus.Desc
	pending_cpus *prometheus.Desc
	running      *prometheus.Desc
	running_cpus *prometheus.Desc
	suspended    *prometheus.Desc
}

func NewUsersCollector(ctx context.Context) *UsersCollector {
	labels := []string{"user"}
	return &UsersCollector{
		ctx:          ctx,
		pending:      prometheus.NewDesc("slurm_user_jobs_pending", "Pending jobs for user", labels, nil),
		pending_cpus: prometheus.NewDesc("slurm_user_cpus_pending", "Pending jobs for user", labels, nil),
		running:      prometheus.NewDesc("slurm_user_jobs_running", "Running jobs for user", labels, nil),
		running_cpus: prometheus.NewDesc("slurm_user_cpus_running", "Running cpus for user", labels, nil),
		suspended:    prometheus.NewDesc("slurm_user_jobs_suspended", "Suspended jobs for user", labels, nil),
	}
}

func (uc *UsersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- uc.pending
	ch <- uc.pending_cpus
	ch <- uc.running
	ch <- uc.running_cpus
	ch <- uc.suspended
}

func (uc *UsersCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := uc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	jobsRespBytes, found := apiCache.Get("jobs")
	if !found {
		slog.Error("failed to get jobs response for users metrics from cache")
		return
	}
	jobsResp, err := api.UnmarshalJobsResponse(jobsRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal jobs response for users metrics", "error", err)
		return
	}
	um, err := ParseUsersMetrics(*jobsResp)
	if err != nil {
		slog.Error("failed to collect user metrics", "error", err)
		return
	}
	for u := range um {
		if um[u].pending > 0 {
			ch <- prometheus.MustNewConstMetric(uc.pending, prometheus.GaugeValue, um[u].pending, u)
		}
		if um[u].pending_cpus > 0 {
			ch <- prometheus.MustNewConstMetric(uc.pending_cpus, prometheus.GaugeValue, um[u].pending_cpus, u)
		}
		if um[u].running > 0 {
			ch <- prometheus.MustNewConstMetric(uc.running, prometheus.GaugeValue, um[u].running, u)
		}
		if um[u].running_cpus > 0 {
			ch <- prometheus.MustNewConstMetric(uc.running_cpus, prometheus.GaugeValue, um[u].running_cpus, u)
		}
		if um[u].suspended > 0 {
			ch <- prometheus.MustNewConstMetric(uc.suspended, prometheus.GaugeValue, um[u].suspended, u)
		}
	}
}

func NewUserJobMetrics() *userJobMetrics {
	return &userJobMetrics{0, 0, 0, 0, 0}
}

type userJobMetrics struct {
	pending      float64
	pending_cpus float64
	running      float64
	running_cpus float64
	suspended    float64
}

func ParseUsersMetrics(jobsResp openapi.V0040OpenapiJobInfoResp) (map[string]*userJobMetrics, error) {
	users := make(map[string]*userJobMetrics)
	for _, j := range jobsResp.Jobs {
		user := *j.UserName
		if _, exists := users[user]; !exists {
			users[user] = NewUserJobMetrics()
		}

		jobState, err := GetJobState(j)
		if err != nil {
			return nil, fmt.Errorf("failed to get job state: %v", err)
		}

		jobCpus, err := GetJobCPUs(j)
		if err != nil {
			return nil, fmt.Errorf("failed to get job cpus: %v", err)
		}

		switch *jobState {
		case types.JobStatePending:
			users[user].pending++
			users[user].pending_cpus += *jobCpus
		case types.JobStateRunning:
			users[user].running++
			users[user].running_cpus += *jobCpus
		case types.JobStateSuspended:
			users[user].suspended++
		}
	}
	return users, nil
}
