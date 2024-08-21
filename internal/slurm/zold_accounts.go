package slurm

import (
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

//
//
// THIS IS ALL OLD DATA FOR CHECKING FEATURE PARITY AND WILL BE REMOVED IN THE FUTURE
//
//
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
