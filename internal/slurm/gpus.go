package slurm

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

type gpusMetrics struct {
	alloc       float64
	idle        float64
	other       float64
	total       float64
	utilization float64
}

func NewGPUsMetrics() *gpusMetrics {
	return &gpusMetrics{}
}

// NOTES:
// node[gres] 		=> gpu:0 										# no gpus
// node[gres] 		=> gpu:nvidia_h100_80gb_hbm3:4(S:0-1) 			# 4 h100 gpus
// node[gres_used]  => gpu:nvidia_h100_80gb_hbm3:4(IDX:0-3) 		# 4 used gpus
// node[gres_used]  => gpu:nvidia_h100_80gb_hbm3:0(IDX:N/A) 		# 0 used gpus
// node[tres]		=> cpu=48,mem=1020522M,billing=48,gres/gpu=4	# 4 total gpus
// node[tres]		=> cpu=1,mem=1M,billing=1						# 0 total gpus
// node[tres_used]	=> cpu=48,mem=1020522M,billing=48,gres/gpu=4	# 4 used gpus
// node[tres_used]	=> cpu=1,mem=1M,billing=1						# 0 used gpus
//
// For tracking gpu resources, it looks like tres will be better. If I need to pull out per-gpu stats later,
// I'll have to use gres
//

// ParseGPUsMetrics iterates through node response objects and tallies up the total and
// allocated gpus, then derives idle and utilization from those numbers.
func ParseGPUsMetrics(nodesResp types.V0040OpenapiNodesResp) (*gpusMetrics, error) {
	gm := NewGPUsMetrics()
	nodes := nodesResp.Nodes
	for _, n := range nodes {
		totalGPUs, err := GetNodeGPUTotal(n)
		if err != nil {
			return nil, fmt.Errorf("failed to get total gpu count for node: %v", err)
		}
		allocGPUs, err := GetNodeGPUAllocated(n)
		if err != nil {
			return nil, fmt.Errorf("failed to get allocated gpu count for node: %v", err)
		}
		idleGPUs := totalGPUs - allocGPUs
		gm.total += float64(totalGPUs)
		gm.alloc += float64(allocGPUs)
		gm.idle += float64(idleGPUs)
	}
	gm.utilization = gm.alloc / gm.total
	return gm, nil
}

/*
 * Implement the Prometheus Collector interface and feed the
 * Slurm scheduler metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */

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

type GPUsCollector struct {
	ctx         context.Context
	alloc       *prometheus.Desc
	idle        *prometheus.Desc
	other       *prometheus.Desc
	total       *prometheus.Desc
	utilization *prometheus.Desc
}

func (cc *GPUsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.other
	ch <- cc.total
	ch <- cc.utilization
}
func (cc *GPUsCollector) Collect(ch chan<- prometheus.Metric) {
	nodeRespBytes, err := api.GetSlurmRestNodesResponse(cc.ctx)
	if err != nil {
		slog.Error("failed to get nodes response for cpu metrics", "error", err)
		return
	}
	nodesResp, err := api.UnmarshalNodesResponse(nodeRespBytes)
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

//
//
// THIS IS ALL OLD DATA FOR CHECKING FEATURE PARITY AND WILL BE REMOVED IN THE FUTURE
//
//

type OldGPUsMetrics struct {
	alloc       float64
	idle        float64
	other       float64
	total       float64
	utilization float64
}

func OldGPUsGetMetrics() *OldGPUsMetrics {
	return OldParseGPUsMetrics()
}

/* TODO:
  sinfo has gresUSED since slurm>=19.05.0rc01 https://github.com/SchedMD/slurm/blob/master/NEWS
  revert to old process on slurm<19.05.0rc01
  --format=AllocGRES will return gres/gpu=8
  --format=AllocTRES will return billing=16,cpu=16,gres/gpu=8,mem=256G,node=1
func ParseAllocatedGPUs() float64 {
	var num_gpus = 0.0

	args := []string{"-a", "-X", "--format=Allocgres", "--state=RUNNING", "--noheader", "--parsable2"}
	output := string(Execute("sacct", args))
	if len(output) > 0 {
		for _, line := range strings.Split(output, "\n") {
			if len(line) > 0 {
				line = strings.Trim(line, "\"")
				descriptor := strings.TrimPrefix(line, "gpu:")
				job_gpus, _ := strconv.ParseFloat(descriptor, 64)
				num_gpus += job_gpus
			}
		}
	}

	return num_gpus
}
*/

func OldParseAllocatedGPUs(data []byte) float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format="Nodes: ,GresUsed:" --state=allocated
	// 3 gpu:2                                       # slurm>=20.11.8
	// 1 gpu:(null):3(IDX:0-7)                       # slurm 21.08.5
	// 13 gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)      # slurm 21.08.5

	sinfo_lines := string(data)
	re := regexp.MustCompile(`gpu:(\(null\)|[^:(]*):?([0-9]+)(\([^)]*\))?`)
	if len(sinfo_lines) > 0 {
		for _, line := range strings.Split(sinfo_lines, "\n") {
			// log.info(line)
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_active_gpus := strings.Fields(line)[1]
				num_node_active_gpus := 0.0
				for _, node_active_gpus_type := range strings.Split(node_active_gpus, ",") {
					if strings.Contains(node_active_gpus_type, "gpu:") {
						node_active_gpus_type = re.FindStringSubmatch(node_active_gpus_type)[2]
						num_node_active_gpus_type, _ := strconv.ParseFloat(node_active_gpus_type, 64)
						num_node_active_gpus += num_node_active_gpus_type
					}
				}
				num_gpus += num_nodes * num_node_active_gpus
			}
		}
	}

	return num_gpus
}

func OldParseIdleGPUs(data []byte) float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format="Nodes: ,Gres: ,GresUsed:" --state=idle,allocated
	// 3 gpu:4 gpu:2                                       																# slurm 20.11.8
	// 1 gpu:8(S:0-1) gpu:(null):3(IDX:0-7)                       												# slurm 21.08.5
	// 13 gpu:A30:4(S:0-1),gpu:Q6K:40(S:0-1) gpu:A30:4(IDX:0-3),gpu:Q6K:4(IDX:0-3)       	# slurm 21.08.5

	sinfo_lines := string(data)
	re := regexp.MustCompile(`gpu:(\(null\)|[^:(]*):?([0-9]+)(\([^)]*\))?`)
	if len(sinfo_lines) > 0 {
		for _, line := range strings.Split(sinfo_lines, "\n") {
			// log.info(line)
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_gpus := strings.Fields(line)[1]
				num_node_gpus := 0.0
				for _, node_gpus_type := range strings.Split(node_gpus, ",") {
					if strings.Contains(node_gpus_type, "gpu:") {
						node_gpus_type = re.FindStringSubmatch(node_gpus_type)[2]
						num_node_gpus_type, _ := strconv.ParseFloat(node_gpus_type, 64)
						num_node_gpus += num_node_gpus_type
					}
				}
				num_node_active_gpus := 0.0
				node_active_gpus := strings.Fields(line)[2]
				for _, node_active_gpus_type := range strings.Split(node_active_gpus, ",") {
					if strings.Contains(node_active_gpus_type, "gpu:") {
						node_active_gpus_type = re.FindStringSubmatch(node_active_gpus_type)[2]
						num_node_active_gpus_type, _ := strconv.ParseFloat(node_active_gpus_type, 64)
						num_node_active_gpus += num_node_active_gpus_type
					}
				}
				num_gpus += num_nodes * (num_node_gpus - num_node_active_gpus)
			}
		}
	}

	return num_gpus
}

func OldParseTotalGPUs(data []byte) float64 {
	var num_gpus = 0.0
	// sinfo -a -h --Format="Nodes: ,Gres:"
	// 3 gpu:4                                       	# slurm 20.11.8
	// 1 gpu:8(S:0-1)                                	# slurm 21.08.5
	// 13 gpu:A30:4(S:0-1),gpu:Q6K:40(S:0-1)        	# slurm 21.08.5

	sinfo_lines := string(data)
	re := regexp.MustCompile(`gpu:(\(null\)|[^:(]*):?([0-9]+)(\([^)]*\))?`)
	if len(sinfo_lines) > 0 {
		for _, line := range strings.Split(sinfo_lines, "\n") {
			// log.Info(line)
			if len(line) > 0 && strings.Contains(line, "gpu:") {
				nodes := strings.Fields(line)[0]
				num_nodes, _ := strconv.ParseFloat(nodes, 64)
				node_gpus := strings.Fields(line)[1]
				num_node_gpus := 0.0
				for _, node_gpus_type := range strings.Split(node_gpus, ",") {
					if strings.Contains(node_gpus_type, "gpu:") {
						node_gpus_type = re.FindStringSubmatch(node_gpus_type)[2]
						num_node_gpus_type, _ := strconv.ParseFloat(node_gpus_type, 64)
						num_node_gpus += num_node_gpus_type
					}
				}
				num_gpus += num_nodes * num_node_gpus
			}
		}
	}

	return num_gpus
}

func OldParseGPUsMetrics() *OldGPUsMetrics {
	var gm OldGPUsMetrics
	total_gpus := OldParseTotalGPUs(OldTotalGPUsData())
	allocated_gpus := OldParseAllocatedGPUs(OldAllocatedGPUsData())
	idle_gpus := OldParseIdleGPUs(OldIdleGPUsData())
	other_gpus := total_gpus - allocated_gpus - idle_gpus
	gm.alloc = allocated_gpus
	gm.idle = idle_gpus
	gm.other = other_gpus
	gm.total = total_gpus
	gm.utilization = allocated_gpus / total_gpus
	return &gm
}

func OldAllocatedGPUsData() []byte {
	args := []string{"-a", "-h", "--Format=Nodes: ,GresUsed:", "--state=allocated"}
	return OldExecute("sinfo", args)
}

func OldIdleGPUsData() []byte {
	args := []string{"-a", "-h", "--Format=Nodes: ,Gres: ,GresUsed:", "--state=idle,allocated"}
	return OldExecute("sinfo", args)
}

func OldTotalGPUsData() []byte {
	args := []string{"-a", "-h", "--Format=Nodes: ,Gres:"}
	return OldExecute("sinfo", args)
}

// OldExecute the sinfo command and return its output
func OldExecute(command string, arguments []string) []byte {
	cmd := exec.Command(command, arguments...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

/*
 * Implement the Prometheus Collector interface and feed the
 * Slurm scheduler metrics into it.
 * https://godoc.org/github.com/prometheus/client_golang/prometheus#Collector
 */

func NewOldGPUsCollector() *OldGPUsCollector {
	return &OldGPUsCollector{
		alloc:       prometheus.NewDesc("slurm_old_gpus_alloc", "Allocated GPUs", nil, nil),
		idle:        prometheus.NewDesc("slurm_old_gpus_idle", "Idle GPUs", nil, nil),
		other:       prometheus.NewDesc("slurm_old_gpus_other", "Other GPUs", nil, nil),
		total:       prometheus.NewDesc("slurm_old_gpus_total", "Total GPUs", nil, nil),
		utilization: prometheus.NewDesc("slurm_old_gpus_utilization", "Total GPU utilization", nil, nil),
	}
}

type OldGPUsCollector struct {
	alloc       *prometheus.Desc
	idle        *prometheus.Desc
	other       *prometheus.Desc
	total       *prometheus.Desc
	utilization *prometheus.Desc
}

// Send all metric descriptions
func (cc *OldGPUsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- cc.alloc
	ch <- cc.idle
	ch <- cc.other
	ch <- cc.total
	ch <- cc.utilization
}
func (cc *OldGPUsCollector) Collect(ch chan<- prometheus.Metric) {
	cm := OldGPUsGetMetrics()
	ch <- prometheus.MustNewConstMetric(cc.alloc, prometheus.GaugeValue, cm.alloc)
	ch <- prometheus.MustNewConstMetric(cc.idle, prometheus.GaugeValue, cm.idle)
	ch <- prometheus.MustNewConstMetric(cc.other, prometheus.GaugeValue, cm.other)
	ch <- prometheus.MustNewConstMetric(cc.total, prometheus.GaugeValue, cm.total)
	ch <- prometheus.MustNewConstMetric(cc.utilization, prometheus.GaugeValue, cm.utilization)
}
