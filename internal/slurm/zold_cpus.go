package slurm

import (
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

//
//
// THIS IS ALL OLD DATA FOR CHECKING FEATURE PARITY AND WILL BE REMOVED IN THE FUTURE
//
//

func CPUsDataOld() []byte {
	cmd := exec.Command("sinfo", "-h", "-o %C")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out, _ := io.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return out
}

func CPUsGetMetricsOld() *cpusMetrics {
	return ParseCPUsMetricsOld(CPUsDataOld())
}

func ParseCPUsMetricsOld(input []byte) *cpusMetrics {
	var cm cpusMetrics
	if strings.Contains(string(input), "/") {
		splitted := strings.Split(strings.TrimSpace(string(input)), "/")
		cm.alloc, _ = strconv.ParseFloat(splitted[0], 64)
		cm.idle, _ = strconv.ParseFloat(splitted[1], 64)
		cm.other, _ = strconv.ParseFloat(splitted[2], 64)
		cm.total, _ = strconv.ParseFloat(splitted[3], 64)
	}
	return &cm
}

func NewCPUsCollectorOld() *CPUsCollectorOld {
	return &CPUsCollectorOld{
		alloc: prometheus.NewDesc("slurm_old_cpus_alloc", "Allocated CPUs", nil, nil),
		idle:  prometheus.NewDesc("slurm_old_cpus_idle", "Idle CPUs", nil, nil),
		other: prometheus.NewDesc("slurm_old_cpus_other", "Mix CPUs", nil, nil),
		total: prometheus.NewDesc("slurm_old_cpus_total", "Total CPUs", nil, nil),
	}
}

type CPUsCollectorOld struct {
	alloc *prometheus.Desc
	idle  *prometheus.Desc
	other *prometheus.Desc
	total *prometheus.Desc
}

// Send all metric descriptions
func (cc *CPUsCollectorOld) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.other
	ch <- cc.total
}
func (cc *CPUsCollectorOld) Collect(ch chan<- prometheus.Metric) {
	cm := CPUsGetMetricsOld()
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
}
