package slurm

import (
	"io"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type OldNodesMetrics struct {
	alloc float64
	comp  float64
	down  float64
	drain float64
	err   float64
	fail  float64
	idle  float64
	maint float64
	mix   float64
	resv  float64
}

func OldNodesGetMetrics() *OldNodesMetrics {
	return OldParseNodesMetrics(OldNodesData())
}

func OldRemoveDuplicates(s []string) []string {
	m := map[string]bool{}
	t := []string{}

	// Walk through the slice 's' and for each value we haven't seen so far, append it to 't'.
	for _, v := range s {
		if _, seen := m[v]; !seen {
			if len(v) > 0 {
				t = append(t, v)
				m[v] = true
			}
		}
	}

	return t
}

func OldParseNodesMetrics(input []byte) *OldNodesMetrics {
	var nm OldNodesMetrics
	lines := strings.Split(string(input), "\n")

	// Sort and remove all the duplicates from the 'sinfo' output
	sort.Strings(lines)
	lines_uniq := OldRemoveDuplicates(lines)

	for _, line := range lines_uniq {
		if strings.Contains(line, ",") {
			split := strings.Split(line, ",")
			count, _ := strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
			state := split[1]
			alloc := regexp.MustCompile(`^alloc`)
			comp := regexp.MustCompile(`^comp`)
			down := regexp.MustCompile(`^down`)
			drain := regexp.MustCompile(`^drain`)
			fail := regexp.MustCompile(`^fail`)
			err := regexp.MustCompile(`^err`)
			idle := regexp.MustCompile(`^idle`)
			maint := regexp.MustCompile(`^maint`)
			mix := regexp.MustCompile(`^mix`)
			resv := regexp.MustCompile(`^res`)
			switch {
			case alloc.MatchString(state) == true:
				nm.alloc += count
			case comp.MatchString(state) == true:
				nm.comp += count
			case down.MatchString(state) == true:
				nm.down += count
			case drain.MatchString(state) == true:
				nm.drain += count
			case fail.MatchString(state) == true:
				nm.fail += count
			case err.MatchString(state) == true:
				nm.err += count
			case idle.MatchString(state) == true:
				nm.idle += count
			case maint.MatchString(state) == true:
				nm.maint += count
			case mix.MatchString(state) == true:
				nm.mix += count
			case resv.MatchString(state) == true:
				nm.resv += count
			}
		}
	}
	return &nm
}

// Execute the sinfo command and return its output
func OldNodesData() []byte {
	cmd := exec.Command("sinfo", "-h", "-o %D,%T")
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

/*
 * Implement the Prometheus Collector interface and feed the
 * Slurm scheduler metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */

func NewOldNodesCollector() *OldNodesCollector {
	return &OldNodesCollector{
		alloc: prometheus.NewDesc("slurm_old_nodes_alloc", "Allocated nodes", nil, nil),
		comp:  prometheus.NewDesc("slurm_old_nodes_comp", "Completing nodes", nil, nil),
		down:  prometheus.NewDesc("slurm_old_nodes_down", "Down nodes", nil, nil),
		drain: prometheus.NewDesc("slurm_old_nodes_drain", "Drain nodes", nil, nil),
		err:   prometheus.NewDesc("slurm_old_nodes_err", "Error nodes", nil, nil),
		fail:  prometheus.NewDesc("slurm_old_nodes_fail", "Fail nodes", nil, nil),
		idle:  prometheus.NewDesc("slurm_old_nodes_idle", "Idle nodes", nil, nil),
		maint: prometheus.NewDesc("slurm_old_nodes_maint", "Maint nodes", nil, nil),
		mix:   prometheus.NewDesc("slurm_old_nodes_mix", "Mix nodes", nil, nil),
		resv:  prometheus.NewDesc("slurm_old_nodes_resv", "Reserved nodes", nil, nil),
	}
}

type OldNodesCollector struct {
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

// Send all metric descriptions
func (nc *OldNodesCollector) Describe(ch chan<- *prometheus.Desc) {
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
func (nc *OldNodesCollector) Collect(ch chan<- prometheus.Metric) {
	nm := OldNodesGetMetrics()
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
