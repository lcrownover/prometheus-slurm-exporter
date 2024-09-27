package slurm

import (
	"context"
	"log/slog"
	"strings"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

/*

AccountsCollector collects metrics for accounts

*/

// AccountsCollector collects metrics for accounts
type AccountsCollector struct {
	ctx          context.Context
	pending      *prometheus.Desc
	pending_cpus *prometheus.Desc
	running      *prometheus.Desc
	running_cpus *prometheus.Desc
	suspended    *prometheus.Desc
}

// NewAccountsCollector creates a new AccountsCollector
func NewAccountsCollector(ctx context.Context) *AccountsCollector {
	labels := []string{"account"}
	return &AccountsCollector{
		ctx:          ctx,
		pending:      prometheus.NewDesc("slurm_account_jobs_pending", "Pending jobs for account", labels, nil),
		pending_cpus: prometheus.NewDesc("slurm_account_cpus_pending", "Pending cpus for account", labels, nil),
		running:      prometheus.NewDesc("slurm_account_jobs_running", "Running jobs for account", labels, nil),
		running_cpus: prometheus.NewDesc("slurm_account_cpus_running", "Running cpus for account", labels, nil),
		suspended:    prometheus.NewDesc("slurm_account_jobs_suspended", "Suspended jobs for account", labels, nil),
	}
}

func (ac *AccountsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ac.pending
	ch <- ac.pending_cpus
	ch <- ac.running
	ch <- ac.running_cpus
	ch <- ac.suspended
}

func (ac *AccountsCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := ac.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	jobsRespBytes, found := apiCache.Get("jobs")
	if !found {
		slog.Error("failed to get jobs response for users metrics from cache")
		return
	}
	jobsResp, err := api.UnmarshalJobsResponse(jobsRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal jobs response for accounts metrics", "error", err)
		return
	}
	am, err := ParseAccountsMetrics(*jobsResp)
	if err != nil {
		slog.Error("failed to parse accounts metrics", "error", err)
		return
	}
	for a := range am {
		if am[a].pending > 0 {
			ch <- prometheus.MustNewConstMetric(ac.pending, prometheus.GaugeValue, am[a].pending, a)
		}
		if am[a].pending_cpus > 0 {
			ch <- prometheus.MustNewConstMetric(ac.pending_cpus, prometheus.GaugeValue, am[a].pending_cpus, a)
		}
		if am[a].running > 0 {
			ch <- prometheus.MustNewConstMetric(ac.running, prometheus.GaugeValue, am[a].running, a)
		}
		if am[a].running_cpus > 0 {
			ch <- prometheus.MustNewConstMetric(ac.running_cpus, prometheus.GaugeValue, am[a].running_cpus, a)
		}
		if am[a].suspended > 0 {
			ch <- prometheus.MustNewConstMetric(ac.suspended, prometheus.GaugeValue, am[a].suspended, a)
		}
	}
}

/*

PartitionsCollector collects metrics for partitions

*/

// CPU metrics collector
type CPUsCollector struct {
	ctx   context.Context
	alloc *prometheus.Desc
	idle  *prometheus.Desc
	other *prometheus.Desc
	total *prometheus.Desc
}

// NewCPUsCollector creates a new CPUsCollector
func NewCPUsCollector(ctx context.Context) *CPUsCollector {
	return &CPUsCollector{
		ctx:   ctx,
		alloc: prometheus.NewDesc("slurm_cpus_alloc", "Allocated CPUs", nil, nil),
		idle:  prometheus.NewDesc("slurm_cpus_idle", "Idle CPUs", nil, nil),
		other: prometheus.NewDesc("slurm_cpus_other", "Mix CPUs", nil, nil),
		total: prometheus.NewDesc("slurm_cpus_total", "Total CPUs", nil, nil),
	}
}

func (cc *CPUsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.other
	ch <- cc.total
}

func (cc *CPUsCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := cc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	jobsRespBytes, found := apiCache.Get("jobs")
	if !found {
		slog.Error("failed to get jobs response for users metrics from cache")
		return
	}
	jobsResp, err := api.UnmarshalJobsResponse(jobsRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal jobs response for cpu metrics", "error", err)
		return
	}
	nodeRespBytes, found := apiCache.Get("nodes")
	if !found {
		slog.Error("failed to get nodes response for cpu metrics from cache")
		return
	}
	nodesResp, err := api.UnmarshalNodesResponse(nodeRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal nodes response for cpu metrics", "error", err)
		return
	}
	cm, err := ParseCPUsMetrics(*nodesResp, *jobsResp)
	if err != nil {
		slog.Error("failed to collect cpus metrics", "error", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
}

type GPUsCollector struct {
	ctx         context.Context
	alloc       *prometheus.Desc
	idle        *prometheus.Desc
	other       *prometheus.Desc
	total       *prometheus.Desc
	utilization *prometheus.Desc
}

func NewGPUsCollector(ctx context.Context) *GPUsCollector {
	return &GPUsCollector{
		ctx:         ctx,
		alloc:       prometheus.NewDesc("slurm_gpus_alloc", "Allocated GPUs", nil, nil),
		idle:        prometheus.NewDesc("slurm_gpus_idle", "Idle GPUs", nil, nil),
		other:       prometheus.NewDesc("slurm_gpus_other", "Other GPUs", nil, nil),
		total:       prometheus.NewDesc("slurm_gpus_total", "Total GPUs", nil, nil),
		utilization: prometheus.NewDesc("slurm_gpus_utilization", "Total GPU utilization", nil, nil),
	}
}

func (cc *GPUsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.other
	ch <- cc.total
	ch <- cc.utilization
}
func (cc *GPUsCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := cc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	nodeRespBytes, found := apiCache.Get("nodes")
	if !found {
		slog.Error("failed to get nodes response for cpu metrics from cache")
		return
	}
	nodesResp, err := api.UnmarshalNodesResponse(nodeRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal nodes response for cpu metrics", "error", err)
		return
	}
	gm, err := ParseGPUsMetrics(*nodesResp)
	if err != nil {
		slog.Error("failed to collect gpus metrics", "error", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, gm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, gm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, gm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, gm.total)
	ch <- prometheus.MustNewConstMetric(cc.utilization, prometheus.GaugeValue, gm.utilization)
}

type NodeCollector struct {
	ctx      context.Context
	cpuAlloc *prometheus.Desc
	cpuIdle  *prometheus.Desc
	cpuOther *prometheus.Desc
	cpuTotal *prometheus.Desc
	memAlloc *prometheus.Desc
	memTotal *prometheus.Desc
}

// NewNodeCollectorOld creates a Prometheus collector to keep all our stats in
// It returns a set of collections for consumption
func NewNodeCollector(ctx context.Context) *NodeCollector {
	labels := []string{"node", "status"}

	return &NodeCollector{
		ctx:      ctx,
		cpuAlloc: prometheus.NewDesc("slurm_node_cpu_alloc", "Allocated CPUs per node", labels, nil),
		cpuIdle:  prometheus.NewDesc("slurm_node_cpu_idle", "Idle CPUs per node", labels, nil),
		cpuOther: prometheus.NewDesc("slurm_node_cpu_other", "Other CPUs per node", labels, nil),
		cpuTotal: prometheus.NewDesc("slurm_node_cpu_total", "Total CPUs per node", labels, nil),
		memAlloc: prometheus.NewDesc("slurm_node_mem_alloc", "Allocated memory per node", labels, nil),
		memTotal: prometheus.NewDesc("slurm_node_mem_total", "Total memory per node", labels, nil),
	}
}

// Send all metric descriptions
func (nc *NodeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.cpuAlloc
	ch <- nc.cpuIdle
	ch <- nc.cpuOther
	ch <- nc.cpuTotal
	ch <- nc.memAlloc
	ch <- nc.memTotal
}

func (nc *NodeCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := nc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	nodeRespBytes, found := apiCache.Get("nodes")
	if !found {
		slog.Error("failed to get nodes response for cpu metrics from cache")
		return
	}
	nodesResp, err := api.UnmarshalNodesResponse(nodeRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal nodes response for cpu metrics", "error", err)
		return
	}
	nm, err := ParseNodeMetrics(*nodesResp)
	if err != nil {
		slog.Error("failed to collect nodes metrics", "error", err)
		return
	}
	for node := range nm {
		ch <- prometheus.MustNewConstMetric(nc.cpuAlloc, prometheus.GaugeValue, float64(nm[node].cpuAlloc), node, nm[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuIdle, prometheus.GaugeValue, float64(nm[node].cpuIdle), node, nm[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuOther, prometheus.GaugeValue, float64(nm[node].cpuOther), node, nm[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.cpuTotal, prometheus.GaugeValue, float64(nm[node].cpuTotal), node, nm[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.memAlloc, prometheus.GaugeValue, float64(nm[node].memAlloc), node, nm[node].nodeStatus)
		ch <- prometheus.MustNewConstMetric(nc.memTotal, prometheus.GaugeValue, float64(nm[node].memTotal), node, nm[node].nodeStatus)
	}
}

type NodesCollector struct {
	ctx   context.Context
	alloc *prometheus.Desc
	comp  *prometheus.Desc
	down  *prometheus.Desc
	drain *prometheus.Desc
	err   *prometheus.Desc
	fail  *prometheus.Desc
	idle  *prometheus.Desc
	maint *prometheus.Desc
	mix   *prometheus.Desc
	resv  *prometheus.Desc
}

func NewNodesCollector(ctx context.Context) *NodesCollector {
	return &NodesCollector{
		ctx:   ctx,
		alloc: prometheus.NewDesc("slurm_nodes_alloc", "Allocated nodes", nil, nil),
		comp:  prometheus.NewDesc("slurm_nodes_comp", "Completing nodes", nil, nil),
		down:  prometheus.NewDesc("slurm_nodes_down", "Down nodes", nil, nil),
		drain: prometheus.NewDesc("slurm_nodes_drain", "Drain nodes", nil, nil),
		err:   prometheus.NewDesc("slurm_nodes_err", "Error nodes", nil, nil),
		fail:  prometheus.NewDesc("slurm_nodes_fail", "Fail nodes", nil, nil),
		idle:  prometheus.NewDesc("slurm_nodes_idle", "Idle nodes", nil, nil),
		maint: prometheus.NewDesc("slurm_nodes_maint", "Maint nodes", nil, nil),
		mix:   prometheus.NewDesc("slurm_nodes_mix", "Mix nodes", nil, nil),
		resv:  prometheus.NewDesc("slurm_nodes_resv", "Reserved nodes", nil, nil),
	}
}

func (nc *NodesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- nc.alloc
	ch <- nc.comp
	ch <- nc.down
	ch <- nc.drain
	ch <- nc.err
	ch <- nc.fail
	ch <- nc.idle
	ch <- nc.maint
	ch <- nc.mix
	ch <- nc.resv
}

func (nc *NodesCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := nc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	nodeRespBytes, found := apiCache.Get("nodes")
	if !found {
		slog.Error("failed to get nodes response for cpu metrics from cache")
		return
	}
	nodesResp, err := api.UnmarshalNodesResponse(nodeRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal nodes response for cpu metrics", "error", err)
		return
	}
	nm, err := ParseNodesMetrics(*nodesResp)
	if err != nil {
		slog.Error("failed to collect nodes metrics", "error", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(nc.alloc, prometheus.GaugeValue, nm.alloc)
	ch <- prometheus.MustNewConstMetric(nc.comp, prometheus.GaugeValue, nm.comp)
	ch <- prometheus.MustNewConstMetric(nc.down, prometheus.GaugeValue, nm.down)
	ch <- prometheus.MustNewConstMetric(nc.drain, prometheus.GaugeValue, nm.drain)
	ch <- prometheus.MustNewConstMetric(nc.err, prometheus.GaugeValue, nm.err)
	ch <- prometheus.MustNewConstMetric(nc.fail, prometheus.GaugeValue, nm.fail)
	ch <- prometheus.MustNewConstMetric(nc.idle, prometheus.GaugeValue, nm.idle)
	ch <- prometheus.MustNewConstMetric(nc.maint, prometheus.GaugeValue, nm.maint)
	ch <- prometheus.MustNewConstMetric(nc.mix, prometheus.GaugeValue, nm.mix)
	ch <- prometheus.MustNewConstMetric(nc.resv, prometheus.GaugeValue, nm.resv)
}

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
	jobsResp, err := api.UnmarshalJobsResponse(jobsRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal jobs response for queue metrics", "error", err)
		return
	}
	qm, err := ParseQueueMetrics(*jobsResp)
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

type SchedulerCollector struct {
	ctx                               context.Context
	threads                           *prometheus.Desc
	queue_size                        *prometheus.Desc
	dbd_queue_size                    *prometheus.Desc
	last_cycle                        *prometheus.Desc
	mean_cycle                        *prometheus.Desc
	cycle_per_minute                  *prometheus.Desc
	backfill_last_cycle               *prometheus.Desc
	backfill_mean_cycle               *prometheus.Desc
	backfill_depth_mean               *prometheus.Desc
	total_backfilled_jobs_since_start *prometheus.Desc
	total_backfilled_jobs_since_cycle *prometheus.Desc
	total_backfilled_heterogeneous    *prometheus.Desc
}

func NewSchedulerCollector(ctx context.Context) *SchedulerCollector {
	return &SchedulerCollector{
		ctx: ctx,
		threads: prometheus.NewDesc(
			"slurm_scheduler_threads",
			"Information provided by the Slurm sdiag command, number of scheduler threads ",
			nil,
			nil),
		queue_size: prometheus.NewDesc(
			"slurm_scheduler_queue_size",
			"Information provided by the Slurm sdiag command, length of the scheduler queue",
			nil,
			nil),
		dbd_queue_size: prometheus.NewDesc(
			"slurm_scheduler_dbd_queue_size",
			"Information provided by the Slurm sdiag command, length of the DBD agent queue",
			nil,
			nil),
		last_cycle: prometheus.NewDesc(
			"slurm_scheduler_last_cycle",
			"Information provided by the Slurm sdiag command, scheduler last cycle time in (microseconds)",
			nil,
			nil),
		mean_cycle: prometheus.NewDesc(
			"slurm_scheduler_mean_cycle",
			"Information provided by the Slurm sdiag command, scheduler mean cycle time in (microseconds)",
			nil,
			nil),
		cycle_per_minute: prometheus.NewDesc(
			"slurm_scheduler_cycle_per_minute",
			"Information provided by the Slurm sdiag command, number scheduler cycles per minute",
			nil,
			nil),
		backfill_last_cycle: prometheus.NewDesc(
			"slurm_scheduler_backfill_last_cycle",
			"Information provided by the Slurm sdiag command, scheduler backfill last cycle time in (microseconds)",
			nil,
			nil),
		backfill_mean_cycle: prometheus.NewDesc(
			"slurm_scheduler_backfill_mean_cycle",
			"Information provided by the Slurm sdiag command, scheduler backfill mean cycle time in (microseconds)",
			nil,
			nil),
		backfill_depth_mean: prometheus.NewDesc(
			"slurm_scheduler_backfill_depth_mean",
			"Information provided by the Slurm sdiag command, scheduler backfill mean depth",
			nil,
			nil),
		total_backfilled_jobs_since_start: prometheus.NewDesc(
			"slurm_scheduler_backfilled_jobs_since_start_total",
			"Information provided by the Slurm sdiag command, number of jobs started thanks to backfilling since last slurm start",
			nil,
			nil),
		total_backfilled_jobs_since_cycle: prometheus.NewDesc(
			"slurm_scheduler_backfilled_jobs_since_cycle_total",
			"Information provided by the Slurm sdiag command, number of jobs started thanks to backfilling since last time stats where reset",
			nil,
			nil),
		total_backfilled_heterogeneous: prometheus.NewDesc(
			"slurm_scheduler_backfilled_heterogeneous_total",
			"Information provided by the Slurm sdiag command, number of heterogeneous job components started thanks to backfilling since last Slurm start",
			nil,
			nil),
	}
}

// Send all metric descriptions
func (c *SchedulerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.threads
	ch <- c.queue_size
	ch <- c.dbd_queue_size
	ch <- c.last_cycle
	ch <- c.mean_cycle
	ch <- c.cycle_per_minute
	ch <- c.backfill_last_cycle
	ch <- c.backfill_mean_cycle
	ch <- c.backfill_depth_mean
	ch <- c.total_backfilled_jobs_since_start
	ch <- c.total_backfilled_jobs_since_cycle
	ch <- c.total_backfilled_heterogeneous
}

// Send the values of all metrics
func (sc *SchedulerCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := sc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	diagRespBytes, found := apiCache.Get("diag")
	if !found {
		slog.Error("failed to get diag response for scheduler metrics from cache")
		return
	}
	diagResp, err := api.UnmarshalDiagResponse(diagRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal diag response for scheduler metrics", "error", err)
		return
	}
	sm, err := ParseSchedulerMetrics(*diagResp)
	if err != nil {
		slog.Error("failed to collect scheduler metrics", "error", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(sc.threads, prometheus.GaugeValue, sm.threads)
	ch <- prometheus.MustNewConstMetric(sc.queue_size, prometheus.GaugeValue, sm.queue_size)
	ch <- prometheus.MustNewConstMetric(sc.dbd_queue_size, prometheus.GaugeValue, sm.dbd_queue_size)
	ch <- prometheus.MustNewConstMetric(sc.last_cycle, prometheus.GaugeValue, sm.last_cycle)
	ch <- prometheus.MustNewConstMetric(sc.mean_cycle, prometheus.GaugeValue, sm.mean_cycle)
	ch <- prometheus.MustNewConstMetric(sc.cycle_per_minute, prometheus.GaugeValue, sm.cycle_per_minute)
	ch <- prometheus.MustNewConstMetric(sc.backfill_last_cycle, prometheus.GaugeValue, sm.backfill_last_cycle)
	ch <- prometheus.MustNewConstMetric(sc.backfill_mean_cycle, prometheus.GaugeValue, sm.backfill_mean_cycle)
	ch <- prometheus.MustNewConstMetric(sc.backfill_depth_mean, prometheus.GaugeValue, sm.backfill_depth_mean)
	ch <- prometheus.MustNewConstMetric(sc.total_backfilled_jobs_since_start, prometheus.GaugeValue, sm.total_backfilled_jobs_since_start)
	ch <- prometheus.MustNewConstMetric(sc.total_backfilled_jobs_since_cycle, prometheus.GaugeValue, sm.total_backfilled_jobs_since_cycle)
	ch <- prometheus.MustNewConstMetric(sc.total_backfilled_heterogeneous, prometheus.GaugeValue, sm.total_backfilled_heterogeneous)
}

type FairShareCollector struct {
	ctx       context.Context
	fairshare *prometheus.Desc
}

func NewFairShareCollector(ctx context.Context) *FairShareCollector {
	labels := []string{"account"}
	return &FairShareCollector{
		ctx:       ctx,
		fairshare: prometheus.NewDesc("slurm_account_fairshare", "FairShare for account", labels, nil),
	}
}

func (fsc *FairShareCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- fsc.fairshare
}

func (fsc *FairShareCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := fsc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	sharesRespBytes, found := apiCache.Get("shares")
	if !found {
		slog.Error("failed to get shares response for fair share metrics from cache")
		return
	}
	// this is disgusting but the response has values of "Infinity" which are
	// not json unmarshal-able, so I manually replace all the "Infinity"s with the correct
	// float64 value that represents Infinity.
	// this will be fixed in v0.0.42
	// https://support.schedmd.com/show_bug.cgi?id=20817
	// 
	// https://github.com/lcrownover/prometheus-slurm-exporter/issues/8
	// also reported that folks are getting "inf" back, so I'll protect for that too
	sharesRespBytes = []byte(strings.ReplaceAll(string(sharesRespBytes.([]byte)), "Infinity", "1.7976931348623157e+308"))
	sharesRespBytes = []byte(strings.ReplaceAll(string(sharesRespBytes.([]byte)), "inf", "1.7976931348623157e+308"))

	sharesResp, err := api.UnmarshalSharesResponse(sharesRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal shares response for fair share metrics", "error", err)
		return
	}
	fsm, err := ParseFairShareMetrics(*sharesResp)
	if err != nil {
		slog.Error("failed to collect fair share metrics", "error", err)
		return
	}
	for f := range fsm {
		ch <- prometheus.MustNewConstMetric(fsc.fairshare, prometheus.GaugeValue, fsm[f].fairshare, f)
	}
}

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

type PartitionsCollector struct {
	ctx       context.Context
	allocated *prometheus.Desc
	idle      *prometheus.Desc
	other     *prometheus.Desc
	pending   *prometheus.Desc
	total     *prometheus.Desc
}

func NewPartitionsCollector(ctx context.Context) *PartitionsCollector {
	labels := []string{"partition"}
	return &PartitionsCollector{
		ctx:       ctx,
		allocated: prometheus.NewDesc("slurm_partition_cpus_allocated", "Allocated CPUs for partition", labels, nil),
		idle:      prometheus.NewDesc("slurm_partition_cpus_idle", "Idle CPUs for partition", labels, nil),
		other:     prometheus.NewDesc("slurm_partition_cpus_other", "Other CPUs for partition", labels, nil),
		pending:   prometheus.NewDesc("slurm_partition_jobs_pending", "Pending jobs for partition", labels, nil),
		total:     prometheus.NewDesc("slurm_partition_cpus_total", "Total CPUs for partition", labels, nil),
	}
}

func (pc *PartitionsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- pc.allocated
	ch <- pc.idle
	ch <- pc.other
	ch <- pc.pending
	ch <- pc.total
}

func (pc *PartitionsCollector) Collect(ch chan<- prometheus.Metric) {
	apiCache := pc.ctx.Value(types.ApiCacheKey).(*cache.Cache)
	partitionRespBytes, found := apiCache.Get("partitions")
	if !found {
		slog.Error("failed to get partitions response for partitions metrics from cache")
		return
	}
	partitionsResp, err := api.UnmarshalPartitionsResponse(partitionRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal partitions response for partitions metrics", "error", err)
		return
	}
	jobsRespBytes, found := apiCache.Get("jobs")
	if !found {
		slog.Error("failed to get jobs response for users metrics from cache")
		return
	}
	jobsResp, err := api.UnmarshalJobsResponse(jobsRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal jobs response for partitions metrics", "error", err)
		return
	}
	nodeRespBytes, found := apiCache.Get("nodes")
	if !found {
		slog.Error("failed to get nodes response for cpu metrics from cache")
		return
	}
	nodesResp, err := api.UnmarshalNodesResponse(nodeRespBytes.([]byte))
	if err != nil {
		slog.Error("failed to unmarshal nodes response for partition metrics", "error", err)
		return
	}
	pm, err := ParsePartitionsMetrics(*partitionsResp, *jobsResp, *nodesResp)
	if err != nil {
		slog.Error("failed to collect partitions metrics", "error", err)
		return
	}
	for p := range pm {
		if pm[p].cpus_allocated > 0 {
			ch <- prometheus.MustNewConstMetric(pc.allocated, prometheus.GaugeValue, pm[p].cpus_allocated, p)
		}
		if pm[p].cpus_idle > 0 {
			ch <- prometheus.MustNewConstMetric(pc.idle, prometheus.GaugeValue, pm[p].cpus_idle, p)
		}
		if pm[p].cpus_other > 0 {
			ch <- prometheus.MustNewConstMetric(pc.other, prometheus.GaugeValue, pm[p].cpus_other, p)
		}
		if pm[p].cpus_total > 0 {
			ch <- prometheus.MustNewConstMetric(pc.total, prometheus.GaugeValue, pm[p].cpus_total, p)
		}
		if pm[p].jobs_pending > 0 {
			ch <- prometheus.MustNewConstMetric(pc.pending, prometheus.GaugeValue, pm[p].jobs_pending, p)
		}
	}
}
