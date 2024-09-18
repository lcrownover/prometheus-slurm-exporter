package types

type JobState string

const (
	JobStatePending     JobState = "pending"
	JobStateCompleted   JobState = "pompleted"
	JobStateFailed      JobState = "failed"
	JobStateOutOfMemory JobState = "out_of_memory"
	JobStateRunning     JobState = "running"
	JobStateSuspended   JobState = "suspended"
	JobStateUnknown     JobState = "unknown"
	JobStateTimeout     JobState = "timeout"
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
