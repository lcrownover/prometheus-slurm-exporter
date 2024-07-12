package slurm

import (
	"fmt"
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

// GetJobState retrieves the first job state string from the JobInfo object in lowercase or returns an error
func GetJobState(job types.V0040JobInfo) (*string, error) {
	// job state should be a list of strings, but the spec is []interface{}
	states := job.JobState
	if states == nil {
		// job state is not found in the job response
		return nil, fmt.Errorf("job state not found in job")
	}
	state := (*states)[0].(string)
	state = strings.ToLower(state)
	return &state, nil
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
