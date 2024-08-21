package slurm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
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
	state := (*states)[0].(string)
	state = strings.ToLower(state)

	completed := regexp.MustCompile(`^completed`)
	pending := regexp.MustCompile(`^pending`)
	failed := regexp.MustCompile(`^failed`)
	running := regexp.MustCompile(`^running`)
	suspended := regexp.MustCompile(`^suspended`)
	out_of_memory := regexp.MustCompile(`^out_of_memory`)
	timeout := regexp.MustCompile(`^timeout`)

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

// GetNodeState returns a NodeState unit or returns an error
func GetNodeState(node types.V0040Node) (*types.NodeState, error) {
	states := node.State
	if states == nil {
		// node state is not found in the node response
		return nil, fmt.Errorf("node state not found in node")
	}
	state := (*states)[0].(string)
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
	default:
		return nil, fmt.Errorf("failed to match cpu state against known states: %v", state)
	}

	return &stateUnit, nil
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
