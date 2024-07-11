/* Copyright 2020-2022 Lucas Crownover, Victor Penso

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package slurm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"io"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
	"github.com/prometheus/client_golang/prometheus"
)

type SlurmJobsResponse struct {
	Jobs []SlurmJobResponse `json:"jobs"`
}

type SlurmJobResourcesResponse struct {
	AllocatedCores float64 `json:"allocated_cores"`
}

type SlurmJobResponse struct {
	Account      string                    `json:"account"`
	JobStates    []string                  `json:"job_state"`
	JobResources SlurmJobResourcesResponse `json:"job_resources"`
}

// AccountsData performs the GET request against the SLURM API and returns
// the response body converted to string.
func AccountsData(ctx context.Context) []byte {
	resp, err := util.NewSlurmGETRequest(ctx, util.ApiAccountsEndpointKey)
	if err != nil {
		slog.Error("failed to perform get request for accounts data", "error", err)
	}
	if resp.StatusCode != 200 {
		slog.Error("received incorrect status code for accounts data")
		slog.Debug("debug", "code", resp.StatusCode, "body", string(resp.Body))
	}
	return resp.Body
}

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

// ParseAccountsMetrics receives the response body of jobs from SLURM and
// parses it into a map of "accountName": *JobMetrics
func ParseAccountsMetrics(jsonResponseBytes []byte) map[string]*JobMetrics {
	accounts := make(map[string]*JobMetrics)
	var sj SlurmJobsResponse
	err := json.Unmarshal(jsonResponseBytes, &sj)
	if err != nil {
		slog.Error("failed to unmarshall job response data", "error", err)
	}
	for _, j := range sj.Jobs {
		account := j.Account
		_, key := accounts[account]
		if !key {
			accounts[account] = NewJobMetrics()
		}
		state := j.JobStates[0]
		state = strings.ToLower(state)
		cpus := j.JobResources.AllocatedCores
		pending := regexp.MustCompile(`^pending`)
		running := regexp.MustCompile(`^running`)
		suspended := regexp.MustCompile(`^suspended`)
		switch {
		case pending.MatchString(state):
			accounts[account].pending++
			accounts[account].pending_cpus += cpus
			slog.Debug("adding pending cpus", "account", fmt.Sprintf("%+v", accounts[account]))
		case running.MatchString(state):
			accounts[account].running++
			accounts[account].running_cpus += cpus
		case suspended.MatchString(state):
			accounts[account].suspended++
		}
	}
	slog.Debug("metrics", "metrics", fmt.Sprintf("%+v", accounts))
	return accounts
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
	ctx = context.WithValue(ctx, util.ApiAccountsEndpointKey, "/slurm/v0.0.40/jobs")
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
	am := ParseAccountsMetrics(AccountsData(ac.ctx))
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

func AccountsDataOld() []byte {
	cmd := exec.Command("squeue", "-a", "-r", "-h", "-o %A|%a|%T|%C")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out, _ := io.ReadAll(io.Reader(stdout))
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return out
}

func ParseAccountsMetricsOld(input []byte) map[string]*JobMetrics {
	accounts := make(map[string]*JobMetrics)
	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		if strings.Contains(line, "|") {
			account := strings.Split(line, "|")[1]
			_, key := accounts[account]
			if !key {
				accounts[account] = &JobMetrics{0, 0, 0, 0, 0}
			}
			state := strings.Split(line, "|")[2]
			state = strings.ToLower(state)
			cpus, _ := strconv.ParseFloat(strings.Split(line, "|")[3], 64)
			pending := regexp.MustCompile(`^pending`)
			running := regexp.MustCompile(`^running`)
			suspended := regexp.MustCompile(`^suspended`)
			switch {
			case pending.MatchString(state):
				accounts[account].pending++
				accounts[account].pending_cpus += cpus
			case running.MatchString(state):
				accounts[account].running++
				accounts[account].running_cpus += cpus
			case suspended.MatchString(state):
				accounts[account].suspended++
			}
		}
	}
	return accounts
}

type OldAccountsCollector struct {
	pending      *prometheus.Desc
	pending_cpus *prometheus.Desc
	running      *prometheus.Desc
	running_cpus *prometheus.Desc
	suspended    *prometheus.Desc
}

func NewOldAccountsCollector() *OldAccountsCollector {
	labels := []string{"account"}
	return &OldAccountsCollector{
		pending:      prometheus.NewDesc("slurm_old_account_jobs_pending", "OLD Pending jobs for account", labels, nil),
		pending_cpus: prometheus.NewDesc("slurm_old_account_cpus_pending", "OLD Pending cpus for account", labels, nil),
		running:      prometheus.NewDesc("slurm_old_account_jobs_running", "OLD Running jobs for account", labels, nil),
		running_cpus: prometheus.NewDesc("slurm_old_account_cpus_running", "OLD Running cpus for account", labels, nil),
		suspended:    prometheus.NewDesc("slurm_old_account_jobs_suspended", "OLD Suspended jobs for account", labels, nil),
	}
}

func (ac *OldAccountsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ac.pending
	ch <- ac.pending_cpus
	ch <- ac.running
	ch <- ac.running_cpus
	ch <- ac.suspended
}

func (ac *OldAccountsCollector) Collect(ch chan<- prometheus.Metric) {
	am := ParseAccountsMetricsOld(AccountsDataOld())
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
