//go:build 2311

package slurm

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

// GetJobAccountName retrieves the account name string from the JobInfo object or returns error
func GetJobAccountName(job types.V0040JobInfo) (*string, error) {
	name := job.Account
	if name == nil {
		return nil, fmt.Errorf("account name not found in job")
	}
	return name, nil
}

// GetJobState returns a JobState unit or returns an error
func GetJobState(job types.V0040JobInfo) (*types.JobState, error) {
	states := job.JobState
	if states == nil {
		// job state is not found in the job response
		return nil, fmt.Errorf("job state not found in job")
	}
	state := string((*states)[0])
	state = strings.ToLower(state)

	completed := regexp.MustCompile(`^completed`)
	pending := regexp.MustCompile(`^pending`)
	failed := regexp.MustCompile(`^failed`)
	running := regexp.MustCompile(`^running`)
	suspended := regexp.MustCompile(`^suspended`)
	out_of_memory := regexp.MustCompile(`^out_of_memory`)
	timeout := regexp.MustCompile(`^timeout`)
	cancelled := regexp.MustCompile(`^cancelled`)
	completing := regexp.MustCompile(`^completing`)
	configuring := regexp.MustCompile(`^configuring`)
	node_fail := regexp.MustCompile(`^node_fail`)
	preempted := regexp.MustCompile(`^preempted`)

	var stateUnit types.JobState

	switch {
	case completed.MatchString(state):
		stateUnit = types.JobStateCompleted
	case pending.MatchString(state):
		stateUnit = types.JobStatePending
	case failed.MatchString(state):
		stateUnit = types.JobStateFailed
	case running.MatchString(state):
		stateUnit = types.JobStateRunning
	case suspended.MatchString(state):
		stateUnit = types.JobStateSuspended
	case out_of_memory.MatchString(state):
		stateUnit = types.JobStateOutOfMemory
	case timeout.MatchString(state):
		stateUnit = types.JobStateTimeout
	case cancelled.MatchString(state):
		stateUnit = types.JobStateCancelled
	case completing.MatchString(state):
		stateUnit = types.JobStateCompleting
	case configuring.MatchString(state):
		stateUnit = types.JobStateConfiguring
	case node_fail.MatchString(state):
		stateUnit = types.JobStateNodeFail
	case preempted.MatchString(state):
		stateUnit = types.JobStatePreempted
	default:
		return nil, fmt.Errorf("failed to match job state against known states: %v", state)
	}

	return &stateUnit, nil
}

// GetJobCPUs retrieves the count of CPUs for the given job or returns an error
func GetJobCPUs(job types.V0040JobInfo) (*float64, error) {
	cn := job.Cpus.Number
	if cn == nil {
		return nil, fmt.Errorf("failed to find cpu count in job")
	}
	cpus := float64(*cn)
	return &cpus, nil
}

// GetNodeStates returns a slice of NodeState unit or returns an error
func GetNodeStates(node types.V0040Node) (*[]types.NodeState, error) {
	var nodeStates []types.NodeState
	states := node.State

	if states == nil {
		// node state is not found in the node response
		return nil, fmt.Errorf("node state not found in node")
	}

	for _, s := range *states {
		state := string(s)
		state = strings.ToLower(state)

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
		notresp := regexp.MustCompile(`^not_responding`)
		invalidreg := regexp.MustCompile(`^invalid_reg`)

		var stateUnit types.NodeState

		switch {
		case alloc.MatchString(state):
			stateUnit = types.NodeStateAlloc
		case comp.MatchString(state):
			stateUnit = types.NodeStateComp
		case down.MatchString(state):
			stateUnit = types.NodeStateDown
		case drain.MatchString(state):
			stateUnit = types.NodeStateDrain
		case fail.MatchString(state):
			stateUnit = types.NodeStateFail
		case err.MatchString(state):
			stateUnit = types.NodeStateErr
		case idle.MatchString(state):
			stateUnit = types.NodeStateIdle
		case maint.MatchString(state):
			stateUnit = types.NodeStateMaint
		case mix.MatchString(state):
			stateUnit = types.NodeStateMix
		case resv.MatchString(state):
			stateUnit = types.NodeStateResv
		case notresp.MatchString(state):
			stateUnit = types.NodeStateNotResponding
		case invalidreg.MatchString(state):
			stateUnit = types.NodeStateInvalidReg
		default:
			return nil, fmt.Errorf("failed to match cpu state against known states: %v", state)
		}

		nodeStates = append(nodeStates, stateUnit)
	}

	return &nodeStates, nil
}

// GetNodeStatesString returns a string of node states separated by delim
func GetNodeStatesString(node types.V0040Node, delim string) (string, error) {
	states, err := GetNodeStates(node)
	if err != nil {
		return "", fmt.Errorf("failed to get node states: %v", err)
	}
	// convert nodestates into strings
	strStates := make([]string, len(*states))
	for i, s := range *states {
		strStates[i] = string(s)
	}
	return strings.Join(strStates, delim), nil
}

// GetGPUTotal returns the number of GPUs in the node
func GetNodeGPUTotal(node types.V0040Node) (int, error) {
	tres := node.Tres
	parts := strings.Split(*tres, ",")
	for _, p := range parts {
		if strings.Contains(p, "gres/gpu=") {
			gp := strings.Split(p, "=")
			if len(gp) != 2 {
				return 0, fmt.Errorf("found gpu in tres but failed to parse: %s", p)
			}
			ns := gp[1]
			n, err := strconv.Atoi(ns)
			if err != nil {
				return 0, fmt.Errorf("failed to parse number of gpus from tres: %s", p)
			}
			return n, nil
		}
	}
	return 0, nil
}

// GetGPUAllocated returns the number of GPUs in the node
func GetNodeGPUAllocated(node types.V0040Node) (int, error) {
	tres := node.TresUsed
	parts := strings.Split(*tres, ",")
	for _, p := range parts {
		if strings.Contains(p, "gres/gpu=") {
			gp := strings.Split(p, "=")
			if len(gp) != 2 {
				return 0, fmt.Errorf("found gpu in tres but failed to parse: %s", p)
			}
			ns := gp[1]
			n, err := strconv.Atoi(ns)
			if err != nil {
				return 0, fmt.Errorf("failed to parse number of gpus from tres: %s", p)
			}
			return n, nil
		}
	}
	return 0, nil
}

// GetNodeAllocMemory returns an unsigned 64bit integer
// of the allocated memory on the node
func GetNodeAllocMemory(node types.V0040Node) uint64 {
	alloc_memory := node.AllocMemory
	return uint64(*alloc_memory)
}

// GetNodeTotalMemory returns an unsigned 64bit integer
// of the total memory on the node
func GetNodeTotalMemory(node types.V0040Node) uint64 {
	total_memory := node.RealMemory
	return uint64(*total_memory)
}

// GetNodeAllocCPUs returns an unsigned 64bit integer
// of the allocated cpus on the node
func GetNodeAllocCPUs(node types.V0040Node) uint64 {
	alloc_cpus := node.AllocCpus
	return uint64(*alloc_cpus)
}

// GetNodeIdleCPUs returns an unsigned 64bit integer
// of the allocated cpus on the node
func GetNodeIdleCPUs(node types.V0040Node) uint64 {
	idle_cpus := node.AllocIdleCpus
	return uint64(*idle_cpus)
}

// GetNodeOtherCPUs returns an unsigned 64bit integer
// of the "other" cpus on the node
// since this isn't in the API, let's just return 0 for now
func GetNodeOtherCPUs(node types.V0040Node) uint64 {
	return 0
}

// GetNodeTotalCPUs returns an unsigned 64bit integer
// of the total cpus on the node
func GetNodeTotalCPUs(node types.V0040Node) uint64 {
	total_cpus := node.Cpus
	return uint64(*total_cpus)
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
		case types.JobStatePending:
			accounts[*account].pending++
			accounts[*account].pending_cpus += *cpus
		case types.JobStateRunning:
			accounts[*account].running++
			accounts[*account].running_cpus += *cpus
		case types.JobStateSuspended:
			accounts[*account].suspended++
		}
	}
	return accounts, nil
}

type cpusMetrics struct {
	alloc float64
	idle  float64
	other float64
	total float64
}

func NewCPUsMetrics() *cpusMetrics {
	return &cpusMetrics{}
}

// ParseCPUMetrics pulls out total cluster cpu states of alloc,idle,other,total
func ParseCPUsMetrics(nodesResp types.V0040OpenapiNodesResp, jobsResp types.V0040OpenapiJobInfoResp) (*cpusMetrics, error) {
	cm := NewCPUsMetrics()
	for _, j := range jobsResp.Jobs {
		state, err := GetJobState(j)
		if err != nil {
			slog.Error("failed to get job state", "error", err)
			continue
		}
		cpus, err := GetJobCPUs(j)
		if err != nil {
			slog.Error("failed to get job cpus", "error", err)
			continue
		}
		// alloc is easy, we just add up all the cpus in the "Running" job state
		if *state == types.JobStateRunning {
			cm.alloc += *cpus
		}
	}
	// total is just the total number of cpus in the cluster
	nodes := nodesResp.Nodes
	for _, n := range nodes {
		if *n.Cpus == 1 {
			// TODO: This probably needs to be a call to partitions to get all nodes
			// in a partition, then add the nodes CPU values up for this field.
			// In our environment, nodes that exist (need slurm commands) get
			// put into slurm without being assigned a partition, but slurm
			// seems to track these systems with cpus=1.
			// This isn't a problem unless your site has nodes with a single CPU.
			continue
		}
		cpus := float64(*n.Cpus)
		cm.total += cpus

		nodeStates, err := GetNodeStates(n)
		if err != nil {
			return nil, fmt.Errorf("failed to get node state for cpu metrics: %v", err)
		}
		for _, ns := range *nodeStates {
			if ns == types.NodeStateMix || ns == types.NodeStateAlloc || ns == types.NodeStateIdle {
				// TODO: This calculate is scuffed. In our 17k core environment, it's
				// reporting ~400 more than the `sinfo -h -o '%C'` command.
				// Gotta figure this one out.
				idle_cpus := float64(*n.AllocIdleCpus)
				cm.idle += idle_cpus
			}
		}
	}
	// Assumedly, this should be fine.
	cm.other = cm.total - cm.idle - cm.alloc
	return cm, nil
}

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
	// TODO: Do we really need an "other" field?
	// using TRES, it should be straightforward.
	gm.utilization = gm.alloc / gm.total
	return gm, nil
}

// NodeMetrics stores metrics for each node
type nodeMetrics struct {
	memAlloc   uint64
	memTotal   uint64
	cpuAlloc   uint64
	cpuIdle    uint64
	cpuOther   uint64
	cpuTotal   uint64
	nodeStatus string
}

func NewNodeMetrics() *nodeMetrics {
	return &nodeMetrics{}
}

// ParseNodeMetrics takes the output of sinfo with node data
// It returns a map of metrics per node
func ParseNodeMetrics(nodesResp types.V0040OpenapiNodesResp) (map[string]*nodeMetrics, error) {
	nodeMap := make(map[string]*nodeMetrics)

	for _, n := range nodesResp.Nodes {
		nodeName := *n.Hostname
		nodeMap[nodeName] = &nodeMetrics{0, 0, 0, 0, 0, 0, ""}

		// state
		nodeStatesStr, err := GetNodeStatesString(n, "|")
		if err != nil {
			return nil, fmt.Errorf("failed to get node state: %v", err)
		}
		nodeMap[nodeName].nodeStatus = nodeStatesStr

		// memory
		nodeMap[nodeName].memAlloc = GetNodeAllocMemory(n)
		nodeMap[nodeName].memTotal = GetNodeTotalMemory(n)

		// cpu
		nodeMap[nodeName].cpuAlloc = GetNodeAllocCPUs(n)
		nodeMap[nodeName].cpuIdle = GetNodeIdleCPUs(n)
		nodeMap[nodeName].cpuOther = GetNodeOtherCPUs(n)
		nodeMap[nodeName].cpuTotal = GetNodeTotalCPUs(n)
	}

	return nodeMap, nil
}

type nodesMetrics struct {
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

func NewNodesMetrics() *nodesMetrics {
	return &nodesMetrics{}
}

// ParseNodesMetrics iterates through node response objects and tallies up
// nodes based on their state
func ParseNodesMetrics(nodesResp types.V0040OpenapiNodesResp) (*nodesMetrics, error) {
	nm := NewNodesMetrics()

	for _, n := range nodesResp.Nodes {
		nodeStates, err := GetNodeStates(n)
		if err != nil {
			return nil, fmt.Errorf("failed to get node state for nodes metrics: %v", err)
		}

		for _, ns := range *nodeStates {
			switch ns {
			case types.NodeStateAlloc:
				nm.alloc += 1
			case types.NodeStateComp:
				nm.comp += 1
			case types.NodeStateDown:
				nm.down += 1
			case types.NodeStateDrain:
				nm.drain += 1
			case types.NodeStateErr:
				nm.err += 1
			case types.NodeStateFail:
				nm.fail += 1
			case types.NodeStateIdle:
				nm.idle += 1
			case types.NodeStateMaint:
				nm.maint += 1
			case types.NodeStateMix:
				nm.mix += 1
			case types.NodeStateResv:
				nm.resv += 1
			}
		}
	}

	return nm, nil
}

func NewQueueMetrics() *queueMetrics {
	return &queueMetrics{}
}

type queueMetrics struct {
	pending     float64
	pending_dep float64
	running     float64
	suspended   float64
	cancelled   float64
	completing  float64
	completed   float64
	configuring float64
	failed      float64
	timeout     float64
	preempted   float64
	node_fail   float64
}

func ParseQueueMetrics(jobsResp types.V0040OpenapiJobInfoResp) (*queueMetrics, error) {
	qm := NewQueueMetrics()
	for _, j := range jobsResp.Jobs {
		jobState, err := GetJobState(j)
		if err != nil {
			return nil, fmt.Errorf("failed to get job state: %v", err)
		}
		switch *jobState {
		case types.JobStatePending:
			if *j.Dependency != "" {
				qm.pending_dep++
			} else {
				qm.pending++
			}
		case types.JobStateRunning:
			qm.running++
		case types.JobStateSuspended:
			qm.suspended++
		case types.JobStateCancelled:
			qm.cancelled++
		case types.JobStateCompleting:
			qm.completing++
		case types.JobStateCompleted:
			qm.completed++
		case types.JobStateConfiguring:
			qm.configuring++
		case types.JobStateFailed:
			qm.failed++
		case types.JobStateTimeout:
			qm.timeout++
		case types.JobStatePreempted:
			qm.preempted++
		case types.JobStateNodeFail:
			qm.node_fail++
		}
	}
	return qm, nil
}

func NewSchedulerMetrics() *schedulerMetrics {
	return &schedulerMetrics{}
}

type schedulerMetrics struct {
	threads                           float64
	queue_size                        float64
	dbd_queue_size                    float64
	last_cycle                        float64
	mean_cycle                        float64
	cycle_per_minute                  float64
	backfill_last_cycle               float64
	backfill_mean_cycle               float64
	backfill_depth_mean               float64
	total_backfilled_jobs_since_start float64
	total_backfilled_jobs_since_cycle float64
	total_backfilled_heterogeneous    float64
}

// Extract the relevant metrics from the sdiag output
func ParseSchedulerMetrics(diagResp types.V0040OpenapiDiagResp) (*schedulerMetrics, error) {
	sm := NewSchedulerMetrics()
	s := diagResp.Statistics

	sm.threads = util.GetValueOrZero(s.ServerThreadCount)
	sm.queue_size = util.GetValueOrZero(s.AgentQueueSize)
	sm.dbd_queue_size = util.GetValueOrZero(s.DbdAgentQueueSize)
	sm.last_cycle = util.GetValueOrZero(s.ScheduleCycleLast)
	sm.mean_cycle = util.GetValueOrZero(s.ScheduleCycleMean)
	sm.cycle_per_minute = util.GetValueOrZero(s.ScheduleCyclePerMinute)
	sm.backfill_depth_mean = util.GetValueOrZero(s.BfDepthMean)
	sm.backfill_last_cycle = util.GetValueOrZero(s.BfCycleLast)
	sm.backfill_mean_cycle = util.GetValueOrZero(s.BfCycleMean)
	sm.total_backfilled_jobs_since_cycle = util.GetValueOrZero(s.BfBackfilledJobs)
	// TODO: This is probably not correct, should revisit this number
	sm.total_backfilled_jobs_since_start = util.GetValueOrZero(s.BfLastBackfilledJobs)
	sm.total_backfilled_heterogeneous = util.GetValueOrZero(s.BfBackfilledHetJobs)
	return sm, nil
}

type fairShareMetrics struct {
	fairshare float64
}

func NewFairShareMetrics() *fairShareMetrics {
	return &fairShareMetrics{}
}

func ParseFairShareMetrics(sharesResp types.V0040OpenapiSharesResp) (map[string]*fairShareMetrics, error) {
	accounts := make(map[string]*fairShareMetrics)
	for _, s := range *sharesResp.Shares.Shares {
		account := *s.Name
		if _, exists := accounts[account]; !exists {
			accounts[account] = NewFairShareMetrics()
		}
		// TODO: check if the level is the right value here,
		// there might be some other property that matches the
		// previous value from the old share info code
		accounts[account].fairshare = *s.Fairshare.Level
	}
	return accounts, nil
}

func NewUserJobMetrics() *userJobMetrics {
	return &userJobMetrics{0, 0, 0, 0, 0}
}

type userJobMetrics struct {
	pending      float64
	pending_cpus float64
	running      float64
	running_cpus float64
	suspended    float64
}

func ParseUsersMetrics(jobsResp types.V0040OpenapiJobInfoResp) (map[string]*userJobMetrics, error) {
	users := make(map[string]*userJobMetrics)
	for _, j := range jobsResp.Jobs {
		user := *j.UserName
		if _, exists := users[user]; !exists {
			users[user] = NewUserJobMetrics()
		}

		jobState, err := GetJobState(j)
		if err != nil {
			return nil, fmt.Errorf("failed to get job state: %v", err)
		}

		jobCpus, err := GetJobCPUs(j)
		if err != nil {
			return nil, fmt.Errorf("failed to get job cpus: %v", err)
		}

		switch *jobState {
		case types.JobStatePending:
			users[user].pending++
			users[user].pending_cpus += *jobCpus
		case types.JobStateRunning:
			users[user].running++
			users[user].running_cpus += *jobCpus
		case types.JobStateSuspended:
			users[user].suspended++
		}
	}
	return users, nil
}
