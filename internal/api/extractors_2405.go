//go:build 2405

package api

import "log/slog"

func ExtractDiagData(diagRespBytes []byte) (*DiagData, error) {
	resp, err := UnmarshalDiagResponse(diagRespBytes)
	if err != nil {
		slog.Error("failed to unmarshal diag response", "error", err)
		return nil, err
	}

	d := &DiagData{
		ApiVersion: "24.05",
	}

	d.SetServerThreadCount(resp.Statistics.ServerThreadCount)
	d.SetDbdAgentQueueSize(resp.Statistics.AgentQueueSize)
	d.SetScheduleCycleLast(resp.Statistics.ScheduleCycleLast)
	d.SetScheduleCycleMean(resp.Statistics.ScheduleCycleMean)
	d.SetScheduleCyclePerMinute(resp.Statistics.ScheduleCyclePerMinute)
	d.SetBfDepthMean(resp.Statistics.BfDepthMean)
	d.SetBfCycleLast(resp.Statistics.BfCycleLast)
	d.SetBfCycleMean(resp.Statistics.BfCycleMean)
	d.SetBfBackfilledJobs(resp.Statistics.BfLastBackfilledJobs)
	d.SetBfBackfilledHetJobs(resp.Statistics.BfBackfilledHetJobs)

	return d, nil
}

func ExtractNodesData(nodesRespBytes []byte) (*NodesData, error) {
	resp, err := UnmarshalNodesResponse(nodesRespBytes)
	if err != nil {
		slog.Error("failed to unmarshal nodes response", "error", err)
		return nil, err
	}

	d := &NodesData{
		ApiVersion: "24.05",
		Nodes:      []NodeData{},
	}

	for _, n := range resp.Nodes {
		nd := NodeData{}
		nd.SetName(n.Name)
		nd.SetHostname(n.Hostname)
		nd.SetNodeStates(n.State)
		nd.SetPartitions(n.Partitions)

		nd.SetTotalCPUs(n.Cpus)
		nd.SetAllocCPUs(n.AllocCpus)
		nd.SetIdleCPUs(n.AllocIdleCpus)
		nd.SetOtherCPUs()

		nd.SetTotalMemory(n.RealMemory)
		nd.SetAllocMemory(n.AllocMemory)

		nd.SetNodeGPUAllocated(n.TresUsed)
		nd.SetNodeGPUTotal(n.Tres)

		d.Nodes = append(d.Nodes, nd)
	}

	return d, nil
}

func ExtractJobsData(jobsRespBytes []byte) (*JobsData, error) {
	resp, err := UnmarshalJobsResponse(jobsRespBytes)
	if err != nil {
		slog.Error("failed to unmarshal jobs response", "error", err)
		return nil, err
	}

	d := &JobsData{
		ApiVersion: "24.05",
		Jobs:       []JobData{},
	}

	for _, j := range resp.Jobs {
		jd := JobData{}
		jd.SetJobAccount(j.Account)
		jd.SetJobPartitionName(j.Partition)
		jd.SetJobState(j.JobState)
		jd.SetJobCPUs(j.JobResources.Cpus)
		jd.SetJobDependency(j.Dependency)

		d.Jobs = append(d.Jobs, jd)
	}

	return d, nil
}

func ExtractPartitionsData(partitionsRespBytes []byte) (*PartitionsData, error) {
	resp, err := UnmarshalPartitionsResponse(partitionsRespBytes)
	if err != nil {
		slog.Error("failed to unmarshal partitions response", "error", err)
		return nil, err
	}

	d := &PartitionsData{
		ApiVersion: "24.05",
		Partitions: []PartitionData{},
	}

	for _, p := range resp.Partitions {
		pd := PartitionData{}
		pd.SetName(p.Name)
		pd.SetTotalCPUs(p.Cpus.Total)
		pd.SetOtherCPUs()
		pd.SetNodeList(p.Nodes.Configured)

		d.Partitions = append(d.Partitions, pd)
	}

	return d, nil
}

func ExtractSharesData(sharesRespBytes []byte) (*SharesData, error) {
	resp, err := UnmarshalSharesResponse(sharesRespBytes)
	if err != nil {
		slog.Error("failed to unmarshal shares response", "error", err)
		return nil, err
	}

	d := &SharesData{
		ApiVersion: "24.05",
		Shares:     []ShareData{},
	}

	for _, s := range resp.Shares.Shares {
		sd := ShareData{}
		sd.SetName(s.Name)
		sd.SetEffectiveUsage(s.EffectiveUsage)

		d.Shares = append(d.Shares, sd)
	}

	return d, nil
}
