//go:build 2411

package api

import (
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

var versionedEndpoints = []endpoint{
	{types.ApiJobsEndpointKey, "jobs", "/slurm/v0.0.42/jobs"},
	{types.ApiNodesEndpointKey, "nodes", "/slurm/v0.0.42/nodes"},
	{types.ApiPartitionsEndpointKey, "partitions", "/slurm/v0.0.42/partitions"},
	{types.ApiDiagEndpointKey, "diag", "/slurm/v0.0.42/diag"},
	{types.ApiSharesEndpointKey, "shares", "/slurm/v0.0.42/shares"},
}
