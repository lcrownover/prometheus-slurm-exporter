package slurm

import (
	"context"
	"log/slog"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

type QueueCollector struct {
	ctx         context.Context
	pending     *prometheus.Desc
	pending_dep *prometheus.Desc
	running     *prometheus.Desc
	suspended   *prometheus.Desc
	cancelled   *prometheus.Desc
	completing  *prometheus.Desc
	completed   *prometheus.Desc
	configuring *prometheus.Desc
	failed      *prometheus.Desc
	timeout     *prometheus.Desc
	preempted   *prometheus.Desc
	node_fail   *prometheus.Desc
}

func NewQueueCollector(ctx context.Context) *QueueCollector {
	return &QueueCollector{
		ctx:         ctx,
		pending:     prometheus.NewDesc("slurm_queue_pending", "Pending jobs in queue", nil, nil),
		pending_dep: prometheus.NewDesc("slurm_queue_pending_dependency", "Pending jobs because of dependency in queue", nil, nil),
		running:     prometheus.NewDesc("slurm_queue_running", "Running jobs in the cluster", nil, nil),
		suspended:   prometheus.NewDesc("slurm_queue_suspended", "Suspended jobs in the cluster", nil, nil),
		cancelled:   prometheus.NewDesc("slurm_queue_cancelled", "Cancelled jobs in the cluster", nil, nil),
		completing:  prometheus.NewDesc("slurm_queue_completing", "Completing jobs in the cluster", nil, nil),
		completed:   prometheus.NewDesc("slurm_queue_completed", "Completed jobs in the cluster", nil, nil),
		configuring: prometheus.NewDesc("slurm_queue_configuring", "Configuring jobs in the cluster", nil, nil),
		failed:      prometheus.NewDesc("slurm_queue_failed", "Number of failed jobs", nil, nil),
		timeout:     prometheus.NewDesc("slurm_queue_timeout", "Jobs stopped by timeout", nil, nil),
		preempted:   prometheus.NewDesc("slurm_queue_preempted", "Number of preempted jobs", nil, nil),
		node_fail:   prometheus.NewDesc("slurm_queue_node_fail", "Number of jobs stopped due to node fail", nil, nil),
	}
}

func (qc *QueueCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- qc.pending
	ch <- qc.pending_dep
	ch <- qc.running
	ch <- qc.suspended
	ch <- qc.cancelled
	ch <- qc.completing
	ch <- qc.completed
	ch <- qc.configuring
	ch <- qc.failed
	ch <- qc.timeout
	ch <- qc.preempted
	ch <- qc.node_fail
}

func (qc *QueueCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := qc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	jobsRespBytes, found := apiCache.Get("jobs")
	if !found {
		slog.Error("failed to get jobs response for users metrics from cache")
		return
	}
	jobsData, err := api.ExtractJobsData(jobsRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to extract jobs data for queue metrics", "error", err)
		return
	}
	qm, err := ParseQueueMetrics(jobsData)
	if err != nil {
		slog.Error("failed to collect queue metrics", "error", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(qc.pending, prometheus.GaugeValue, qm.pending)
	ch <- prometheus.MustNewConstMetric(qc.pending_dep, prometheus.GaugeValue, qm.pending_dep)
	ch <- prometheus.MustNewConstMetric(qc.running, prometheus.GaugeValue, qm.running)
	ch <- prometheus.MustNewConstMetric(qc.suspended, prometheus.GaugeValue, qm.suspended)
	ch <- prometheus.MustNewConstMetric(qc.cancelled, prometheus.GaugeValue, qm.cancelled)
	ch <- prometheus.MustNewConstMetric(qc.completing, prometheus.GaugeValue, qm.completing)
	ch <- prometheus.MustNewConstMetric(qc.completed, prometheus.GaugeValue, qm.completed)
	ch <- prometheus.MustNewConstMetric(qc.configuring, prometheus.GaugeValue, qm.configuring)
	ch <- prometheus.MustNewConstMetric(qc.failed, prometheus.GaugeValue, qm.failed)
	ch <- prometheus.MustNewConstMetric(qc.timeout, prometheus.GaugeValue, qm.timeout)
	ch <- prometheus.MustNewConstMetric(qc.preempted, prometheus.GaugeValue, qm.preempted)
	ch <- prometheus.MustNewConstMetric(qc.node_fail, prometheus.GaugeValue, qm.node_fail)
}

func NewQueueMetrics() *queueMetrics {
	return &queueMetrics{}
}

type queueMetrics struct {
	pending     float64
	pending_dep float64
	running     float64
	suspended   float64
	cancelled   float64
	completing  float64
	completed   float64
	configuring float64
	failed      float64
	timeout     float64
	preempted   float64
	node_fail   float64
}

func ParseQueueMetrics(jobsData *api.JobsData) (*queueMetrics, error) {
	qm := NewQueueMetrics()
	for _, j := range jobsData.Jobs {
		switch j.JobState {
		case types.JobStatePending:
			if j.Dependency != "" {
				qm.pending_dep++
			} else {
				qm.pending++
			}
		case types.JobStateRunning:
			qm.running++
		case types.JobStateSuspended:
			qm.suspended++
		case types.JobStateCancelled:
			qm.cancelled++
		case types.JobStateCompleting:
			qm.completing++
		case types.JobStateCompleted:
			qm.completed++
		case types.JobStateConfiguring:
			qm.configuring++
		case types.JobStateFailed:
			qm.failed++
		case types.JobStateTimeout:
			qm.timeout++
		case types.JobStatePreempted:
			qm.preempted++
		case types.JobStateNodeFail:
			qm.node_fail++
		}
	}
	return qm, nil
}
