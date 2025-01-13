//go:build 2411

package api

var apiVersion = "24.11"

type DiagResp struct {
	Statistics struct {
		ServerThreadCount      *int32 `json:"server_thread_count"`
		AgentQueueSize         *int32 `json:"agent_queue_size"`
		DbdAgentQueueSize      *int32 `json:"dbd_agent_queue_size"`
		ScheduleCycleLast      *int32 `json:"schedule_cycle_last"`
		ScheduleCycleMean      *int64 `json:"schedule_cycle_mean"`
		ScheduleCyclePerMinute *int64 `json:"schedule_cycle_per_minute"`
		BfDepthMean            *int64 `json:"bf_depth_mean"`
		BfCycleLast            *int32 `json:"bf_cycle_last"`
		BfCycleMean            *int64 `json:"bf_cycle_mean"`
		BfBackfilledJobs       *int32 `json:"bf_backfilled_jobs"`
		BfLastBackfilledJobs   *int32 `json:"bf_last_backfilled_jobs"`
		BfBackfilledHetJobs    *int32 `json:"bf_backfilled_het_jobs"`
	} `json:"statistics"`
}

type JobsResp struct {
	Jobs []struct {
		Account      *string  `json:"account"`
		UserName     *string  `json:"user_name"`
		Partition    *string  `json:"partition"`
		JobState     []string `json:"job_state"`
		Dependency   *string  `json:"dependency"`
		JobResources struct {
			Cpus *int32 `json:"cpus"`
		} `json:"job_resources"`
	} `json:"jobs"`
}

type NodesResp struct {
	Nodes []struct {
		Name          *string  `json:"name,omitempty"`
		Hostname      *string  `json:"hostname,omitempty"`
		State         []string `json:"state,omitempty"`
		Tres          *string  `json:"tres,omitempty"`
		TresUsed      *string  `json:"tres_used,omitempty"`
		Partitions    []string `json:"partitions,omitempty"`
		AllocMemory   *int64   `json:"alloc_memory,omitempty"`
		RealMemory    *int64   `json:"real_memory,omitempty"`
		AllocCpus     *int32   `json:"alloc_cpus,omitempty"`
		AllocIdleCpus *int32   `json:"alloc_idle_cpus,omitempty"`
		Cpus          *int32   `json:"cpus,omitempty"`
	} `json:"nodes"`
}

type PartitionsResp struct {
	Partitions []struct {
		Name *string `json:"name,omitempty"`
		Cpus *struct {
			Total *int32 `json:"total"`
		} `json:"cpus"`
		Nodes *struct {
			Configured *string `json:"configured"`
		} `json:"nodes"`
	} `json:"partitions"`
}

type SharesResp struct {
	Shares struct {
		Shares []struct {
			Name           *string  `json:"name"`
			EffectiveUsage *struct {
				Number   *float64 `json:"number"`
			} `json:"effective_usage"`
		} `json:"shares"`
	} `json:"shares"`
}
