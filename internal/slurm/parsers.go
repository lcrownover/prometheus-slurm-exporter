package slurm

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

type JobState int

const (
	JobStatePending JobState = iota
	JobStateCompleted
	JobStateFailed
	JobStateOutOfMemory
	JobStateRunning
	JobStateSuspended
	JobStateUnknown
)

type NodeState int

const (
	NodeStateAlloc NodeState = iota
	NodeStateComp
	NodeStateDown
	NodeStateDrain
	NodeStateFail
	NodeStateErr
	NodeStateIdle
	NodeStateMaint
	NodeStateMix
	NodeStateResv
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
func GetJobState(job types.V0040JobInfo) (*JobState, error) {
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

	var stateUnit JobState

	switch {
	case completed.MatchString(state):
		stateUnit = JobStateCompleted
	case pending.MatchString(state):
		stateUnit = JobStatePending
	case failed.MatchString(state):
		stateUnit = JobStateFailed
	case running.MatchString(state):
		stateUnit = JobStateRunning
	case suspended.MatchString(state):
		stateUnit = JobStateSuspended
	case out_of_memory.MatchString(state):
		stateUnit = JobStateOutOfMemory
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
func GetNodeState(node types.V0040Node) (*NodeState, error) {
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

	var stateUnit NodeState

	switch {
	case alloc.MatchString(state):
		stateUnit = NodeStateAlloc
	case comp.MatchString(state):
		stateUnit = NodeStateComp
	case down.MatchString(state):
		stateUnit = NodeStateDown
	case drain.MatchString(state):
		stateUnit = NodeStateDrain
	case fail.MatchString(state):
		stateUnit = NodeStateFail
	case err.MatchString(state):
		stateUnit = NodeStateErr
	case idle.MatchString(state):
		stateUnit = NodeStateIdle
	case maint.MatchString(state):
		stateUnit = NodeStateMaint
	case mix.MatchString(state):
		stateUnit = NodeStateMix
	case resv.MatchString(state):
		stateUnit = NodeStateResv
	default:
		return nil, fmt.Errorf("failed to match cpu state against known states: %v", state)
	}

	return &stateUnit, nil
}
