package slurm

import (
	"context"
	"log/slog"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

type JobMetrics struct {
	pending      float64
	pending_cpus float64
	running      float64
	running_cpus float64
	suspended    float64
}

func NewJobMetrics() *JobMetrics {
	return &JobMetrics{}
}

// ParseAccountsMetrics gets the response body of jobs from SLURM and
// parses it into a map of "accountName": *JobMetrics
func ParseAccountsMetrics(jobs []types.V0040JobInfo) (map[string]*JobMetrics, error) {
	accounts := make(map[string]*JobMetrics)
	for _, j := range jobs {
		// get the account name
		account, err := GetJobAccountName(j)
		if err != nil {
			slog.Error("failed to find account name in job", "error", err)
			continue
		}
		// build the map with the account name as the key and job metrics as the value
		_, key := accounts[*account]
		if !key {
			// initialize a new metrics object if the key isnt found
			accounts[*account] = NewJobMetrics()
		}
		// get the job state
		state, err := GetJobState(j)
		if err != nil {
			slog.Error("failed to parse job state", "error", err)
			continue
		}
		// get the cpus for the job
		cpus, err := GetJobCPUs(j)
		if err != nil {
			slog.Error("failed to parse job cpus", "error", err)
			continue
		}
		// for each of the jobs, depending on the state,
		// tally up the cpu count and increment the count of jobs for that state
		switch *state {
		case JobStatePending:
			accounts[*account].pending++
			accounts[*account].pending_cpus += *cpus
		case JobStateRunning:
			accounts[*account].running++
			accounts[*account].running_cpus += *cpus
		case JobStateSuspended:
			accounts[*account].suspended++
		}
	}
	return accounts, nil
}

type AccountsCollector struct {
	ctx          context.Context
	pending      *prometheus.Desc
	pending_cpus *prometheus.Desc
	running      *prometheus.Desc
	running_cpus *prometheus.Desc
	suspended    *prometheus.Desc
}

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
	resp, err := api.GetSlurmRestJobsResponse(ac.ctx)
	if err != nil {
		slog.Error("failed to get jobs response for accounts metrics", "error", err)
		return
	}
	jobResp, err := api.UnmarshalJobsResponse(resp)
	if err != nil {
		slog.Error("failed to unmarshal jobs response for accounts metrics", "error", err)
		return
	}
	am, err := ParseAccountsMetrics(jobResp.Jobs)
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
