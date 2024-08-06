package slurm
//
// import (
//   "context"
//   "log/slog"
// 	"io/ioutil"
// 	"log"
// 	"os/exec"
// 	"regexp"
// 	"strconv"
// 	"strings"
//
// 	"github.com/prometheus/client_golang/prometheus"
//   "github.com/lcrownover/prometheus-slurm-exporter/internal/types"
// )
//
//
//
//
// type UserMetrics struct {
// 	pending      float64
// 	pending_cpus float64
// 	running      float64
// 	running_cpus float64
// 	suspended    float64
// }
//
// func NewJobMetrics() *AccountsJobMetrics {
//   return &AccountsJobMetrics{}
// }
//
// func ParseUserMetrics(jobs []types.V0040JobInfo) (map[string]*AccountsJobMetrics, error) {
// 	users := make(map[string]*AccountsJobMetrics)
// 	for _, j := range jobs {
// 		userName, err := GetJobAccountName(j)
// 		if err != nil {
// 			slog.Error("failed to find user name in job", "error", err)
// 			continue
// 		}
// 		if _, key := users[*userName]; !key {
// 			users[*userName] = NewJobMetrics()
// 		}
//
// 		state, err := GetJobState(j)
// 		if err != nil {
// 			slog.Error("failed to parse job state", "error", err)
// 			continue
// 		}
//
// 		cpus, err := GetJobCPUs(j)
// 		if err != nil {
// 			slog.Error("failed to parse job cpus", "error", err)
// 			continue
// 		}
// 		switch *state {
// 		case JobStatePending:
// 			users[*userName].pending++
// 			users[*userName].pending_cpus += *cpus
// 		case JobStateRunning:
// 			users[*userName].running++
// 			users[*userName].running_cpus += *cpus
// 		case JobStateSuspended:
// 			users[*userName].suspended++
// 		}
// 	}
// 	return users, nil
// }
//
//
// type UsersCollector struct {
//   ctx           context.Context
// 	pending      *prometheus.Desc
// 	pending_cpus *prometheus.Desc
// 	running      *prometheus.Desc
// 	running_cpus *prometheus.Desc
// 	suspended    *prometheus.Desc
// }
//
//
// func NewUsersCollector() *UsersCollector {
// 	labels := []string{"user"}
// 	return &UsersCollector{
// 		pending:      prometheus.NewDesc("slurm_user_jobs_pending", "Pending jobs for user", labels, nil),
// 		pending_cpus: prometheus.NewDesc("slurm_user_cpus_pending", "Pending jobs for user", labels, nil),
// 		running:      prometheus.NewDesc("slurm_user_jobs_running", "Running jobs for user", labels, nil),
// 		running_cpus: prometheus.NewDesc("slurm_user_cpus_running", "Running cpus for user", labels, nil),
// 		suspended:    prometheus.NewDesc("slurm_user_jobs_suspended", "Suspended jobs for user", labels, nil),
// 	}
// }
//
// func (uc *UsersCollector) Describe(ch chan<- *prometheus.Desc) {
// 	ch <- uc.pending
// 	ch <- uc.pending_cpus
// 	ch <- uc.running
// 	ch <- uc.running_cpus
// 	ch <- uc.suspended
// }
//
// func (uc *UsersCollector) Collect(ch chan<- prometheus.Metric) {
// 	resp, err := getSlurmRestJobsResponse()
// 	if err != nil {
// 		slog.Error("failed to get jobs response for user metrics", "error", err)
// 		return
// 	}
//
// 	var dataMap map[string]interface{}
// 	err = json.Unmarshal(resp, &dataMap)
// 	if err != nil {
// 		slog.Error("failed to unmarshal jobs response for user metrics", "error", err)
// 		return
// 	}
//
// 	jobs, ok := dataMap["jobs"].([]interface{})
// 	if !ok {
// 		slog.Error("\"jobs\" key not found or not an array")
// 		return
// 	}
//
// 	um, err := ParseUserMetrics(jobs)
// 	if err != nil {
// 		slog.Error("failed to parse user metrics", "error", err)
// 		return
// 	}
//
// 	for u := range um {
//
// 		if um[u].pending > 0 {
// 			ch <- prometheus.MustNewConstMetric(uc.pending, prometheus.GaugeValue, um[u].pending, u)
// 		}
// 		if um[u].pending_cpus > 0 {
// 			ch <- prometheus.MustNewConstMetric(uc.pending_cpus, prometheus.GaugeValue, um[u].pending_cpus, u)
// 		}
// 		if um[u].running > 0 {
// 			ch <- prometheus.MustNewConstMetric(uc.running, prometheus.GaugeValue, um[u].running, u)
// 		}
// 		if um[u].running_cpus > 0 {
// 			ch <- prometheus.MustNewConstMetric(uc.running_cpus, prometheus.GaugeValue, um[u].running_cpus, u)
// 		}
// 		if um[u].suspended > 0 {
// 			ch <- prometheus.MustNewConstMetric(uc.suspended, prometheus.GaugeValue, um[u].suspended, u)
// 		}
// 	}
// }
//
//
