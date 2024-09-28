//go:build 2311

package slurm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	openapi "github.com/lcrownover/openapi-slurm-23-11"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

// GetJobAccountName retrieves the account name string from the JobInfo object or returns error
func GetJobAccountName(job openapi.V0040JobInfo) (*string, error) {
	name := job.Account
	if name == nil {
		return nil, fmt.Errorf("account name not found in job")
	}
	return name, nil
}

// GetJobPartitionName retrieves the partition name string from the JobInfo object or returns error
func GetJobPartitionName(job openapi.V0040JobInfo) (*string, error) {
	name := job.Partition
	if name == nil {
		return nil, fmt.Errorf("partition name not found in job")
	}
	return name, nil
}

// GetJobState returns a JobState unit or returns an error
func GetJobState(job openapi.V0040JobInfo) (*types.JobState, error) {
	states := job.JobState
	if states == nil {
		// job state is not found in the job response
		return nil, fmt.Errorf("job state not found in job")
	}
	state := string((states)[0])
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

// GetNodeName retrieves the node name string from the Node object or returns error
func GetNodeName(node openapi.V0040Node) (*string, error) {
	name := node.Name
	if name == nil {
		return nil, fmt.Errorf("node name not found in node information")
	}
	return name, nil
}

// GetJobCPUs retrieves the count of CPUs for the given job or returns an error
func GetJobCPUs(job openapi.V0040JobInfo) (*float64, error) {
	cn := job.Cpus.Number
	if cn == nil {
		return nil, fmt.Errorf("failed to find cpu count in job")
	}
	cpus := float64(*cn)
	return &cpus, nil
}

// GetPartitionName retrieves the name for a given partition or returns an error
func GetPartitionName(partition openapi.V0040PartitionInfo) (*string, error) {
	pn := partition.Name
	if pn == nil {
		return nil, fmt.Errorf("failed to find name in partition")
	}
	return pn, nil
}

// GetPartitionTotalCPUs retrieves the count of total CPUs for a given partition or returns an error
func GetPartitionTotalCPUs(partition openapi.V0040PartitionInfo) (*float64, error) {
	pn := partition.Cpus.Total
	if pn == nil {
		return nil, fmt.Errorf("failed to find total cpus in partition")
	}
	cpus := float64(*pn)
	return &cpus, nil
}

// GetPartitionNodeList retrieves the slurm node notation for nodes assigned to the partition
// returns an empty string if none found, only errors on nil pointer from json
func GetPartitionNodeList(partition openapi.V0040PartitionInfo) (string, error) {
	nodeList := partition.Nodes.Configured
	if nodeList == nil {
		return "", fmt.Errorf("failed to find total cpus in partition")
	}
	return *nodeList, nil
}

// GetNodeStates returns a slice of NodeState unit or returns an error
func GetNodeStates(node openapi.V0040Node) (*[]types.NodeState, error) {
	var nodeStates []types.NodeState
	states := node.State

	if states == nil {
		// node state is not found in the node response
		return nil, fmt.Errorf("node state not found in node")
	}

	for _, s := range states {
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
		planned := regexp.MustCompile(`^planned`)
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
		case planned.MatchString(state):
			stateUnit = types.NodeStatePlanned
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
func GetNodeStatesString(node openapi.V0040Node, delim string) (string, error) {
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
func GetNodeGPUTotal(node openapi.V0040Node) (int, error) {
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
func GetNodeGPUAllocated(node openapi.V0040Node) (int, error) {
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

// GetNodePartitions returns a list of strings that are the partitions a node belongs to
func GetNodePartitions(node openapi.V0040Node) []string {
	ps := node.Partitions
	if ps == nil {
		return []string{}
	}
	return ps
}

// GetNodeAllocMemory returns an unsigned 64bit integer
// of the allocated memory on the node
func GetNodeAllocMemory(node openapi.V0040Node) uint64 {
	alloc_memory := node.AllocMemory
	return uint64(*alloc_memory)
}

// GetNodeTotalMemory returns an unsigned 64bit integer
// of the total memory on the node
func GetNodeTotalMemory(node openapi.V0040Node) uint64 {
	total_memory := node.RealMemory
	return uint64(*total_memory)
}

// GetNodeAllocCPUs returns an unsigned 64bit integer
// of the allocated cpus on the node
func GetNodeAllocCPUs(node openapi.V0040Node) uint64 {
	alloc_cpus := node.AllocCpus
	return uint64(*alloc_cpus)
}

// GetNodeIdleCPUs returns an unsigned 64bit integer
// of the allocated cpus on the node
func GetNodeIdleCPUs(node openapi.V0040Node) uint64 {
	idle_cpus := node.AllocIdleCpus
	return uint64(*idle_cpus)
}

// GetNodeOtherCPUs returns an unsigned 64bit integer
// of the "other" cpus on the node
// since this isn't in the API, let's just return 0 for now
func GetNodeOtherCPUs(node openapi.V0040Node) uint64 {
	return 0
}

// GetNodeTotalCPUs returns an unsigned 64bit integer
// of the total cpus on the node
func GetNodeTotalCPUs(node openapi.V0040Node) uint64 {
	total_cpus := node.Cpus
	return uint64(*total_cpus)
}

// GetFairShareEffectiveUsage returns a float64
// of the effective usage of the fair share system
func GetFairShareEffectiveUsage(share openapi.V0040AssocSharesObjWrap) float64 {
	eu := share.EffectiveUsage
	if eu == nil {
		return 0
	}
	return *eu
}
