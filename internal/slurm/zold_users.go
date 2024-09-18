package slurm

import (
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func UsersDataOld() []byte {
	cmd := exec.Command("squeue", "-a", "-r", "-h", "-o %A|%u|%T|%C")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out, _ := ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return out
}

type UserJobMetricsOld struct {
	pending      float64
	pending_cpus float64
	running      float64
	running_cpus float64
	suspended    float64
}

func ParseUsersMetricsOld(input []byte) map[string]*UserJobMetricsOld {
	users := make(map[string]*UserJobMetricsOld)
	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		if strings.Contains(line, "|") {
			user := strings.Split(line, "|")[1]
			_, key := users[user]
			if !key {
				users[user] = &UserJobMetricsOld{0, 0, 0, 0, 0}
			}
			state := strings.Split(line, "|")[2]
			state = strings.ToLower(state)
			cpus, _ := strconv.ParseFloat(strings.Split(line, "|")[3], 64)
			pending := regexp.MustCompile(`^pending`)
			running := regexp.MustCompile(`^running`)
			suspended := regexp.MustCompile(`^suspended`)
			switch {
			case pending.MatchString(state) == true:
				users[user].pending++
				users[user].pending_cpus += cpus
			case running.MatchString(state) == true:
				users[user].running++
				users[user].running_cpus += cpus
			case suspended.MatchString(state) == true:
				users[user].suspended++
			}
		}
	}
	return users
}

type UsersCollectorOld struct {
	pending      *prometheus.Desc
	pending_cpus *prometheus.Desc
	running      *prometheus.Desc
	running_cpus *prometheus.Desc
	suspended    *prometheus.Desc
}

func NewUsersCollectorOld() *UsersCollector {
	labels := []string{"user"}
	return &UsersCollector{
		pending:      prometheus.NewDesc("slurm_old_user_jobs_pending", "Pending jobs for user", labels, nil),
		pending_cpus: prometheus.NewDesc("slurm_old_user_cpus_pending", "Pending jobs for user", labels, nil),
		running:      prometheus.NewDesc("slurm_old_user_jobs_running", "Running jobs for user", labels, nil),
		running_cpus: prometheus.NewDesc("slurm_old_user_cpus_running", "Running cpus for user", labels, nil),
		suspended:    prometheus.NewDesc("slurm_old_user_jobs_suspended", "Suspended jobs for user", labels, nil),
	}
}

func (uc *UsersCollectorOld) Describe(ch chan<- *prometheus.Desc) {
	ch <- uc.pending
	ch <- uc.pending_cpus
	ch <- uc.running
	ch <- uc.running_cpus
	ch <- uc.suspended
}

func (uc *UsersCollectorOld) Collect(ch chan<- prometheus.Metric) {
	um := ParseUsersMetricsOld(UsersDataOld())
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
