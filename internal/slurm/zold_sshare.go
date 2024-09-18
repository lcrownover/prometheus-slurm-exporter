package slurm

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func FairShareDataOld() []byte {
	cmd := exec.Command("sshare", "-n", "-P", "-o", "account,fairshare")
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

type FairShareMetricsOld struct {
	fairshare float64
}

func ParseFairShareMetricsOld() map[string]*FairShareMetricsOld {
	accounts := make(map[string]*FairShareMetricsOld)
	lines := strings.Split(string(FairShareDataOld()), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "  ") {
			if strings.Contains(line, "|") {
				account := strings.Trim(strings.Split(line, "|")[0], " ")
				_, key := accounts[account]
				if !key {
					accounts[account] = &FairShareMetricsOld{0}
				}
				fairshare, _ := strconv.ParseFloat(strings.Split(line, "|")[1], 64)
				accounts[account].fairshare = fairshare
			}
		}
	}
	return accounts
}

type FairShareCollectorOld struct {
	fairshare *prometheus.Desc
}

func NewFairShareCollectorOld() *FairShareCollectorOld {
	labels := []string{"account"}
	return &FairShareCollectorOld{
		fairshare: prometheus.NewDesc("slurm_old_account_fairshare", "FairShare for account", labels, nil),
	}
}

func (fsc *FairShareCollectorOld) Describe(ch chan<- *prometheus.Desc) {
	ch <- fsc.fairshare
}

func (fsc *FairShareCollectorOld) Collect(ch chan<- prometheus.Metric) {
	fsm := ParseFairShareMetricsOld()
	for f := range fsm {
		ch <- prometheus.MustNewConstMetric(fsc.fairshare, prometheus.GaugeValue, fsm[f].fairshare, f)
	}
}
