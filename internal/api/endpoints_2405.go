//go:build 2405

package api

import (
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

var versionedEndpoints = []endpoint{
	{types.ApiJobsEndpointKey, "jobs", "/slurm/v0.0.41/jobs"},
	{types.ApiNodesEndpointKey, "nodes", "/slurm/v0.0.41/nodes"},
	{types.ApiPartitionsEndpointKey, "partitions", "/slurm/v0.0.41/partitions"},
	{types.ApiDiagEndpointKey, "diag", "/slurm/v0.0.41/diag"},
	{types.ApiSharesEndpointKey, "shares", "/slurm/v0.0.41/shares"},
}
