//go:build 2311

package api

import (
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

var versionedEndpoints = []endpoint{
	{types.ApiJobsEndpointKey, "jobs", "/slurm/v0.0.40/jobs"},
	{types.ApiNodesEndpointKey, "nodes", "/slurm/v0.0.40/nodes"},
	{types.ApiPartitionsEndpointKey, "partitions", "/slurm/v0.0.40/partitions"},
	{types.ApiDiagEndpointKey, "diag", "/slurm/v0.0.40/diag"},
	{types.ApiSharesEndpointKey, "shares", "/slurm/v0.0.40/shares"},
}
