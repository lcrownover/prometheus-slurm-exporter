package types

type JobState int

const (
	JobStatePending JobState = iota
	JobStateCompleted
	JobStateFailed
	JobStateOutOfMemory
	JobStateRunning
	JobStateSuspended
	JobStateUnknown
	JobStateTimeout
)

type SlurmJobsResponse struct {
	Jobs []slurmJob `json:"jobs"`
}

type slurmJobCPUs struct {
	Number int `json:"number"`
}

type slurmJob struct {
	Account   string       `json:"account"`
	JobStates []string     `json:"job_state"`
	CPUs      slurmJobCPUs `json:"cpus"`
}
