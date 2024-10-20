//go:build 2311

package api

import (
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

var versionedEndpoints = []endpoint{
	{types.ApiJobsEndpointKey, "/slurm/v0.0.40/jobs"},
	{types.ApiNodesEndpointKey, "/slurm/v0.0.40/nodes"},
	{types.ApiPartitionsEndpointKey, "/slurm/v0.0.40/partitions"},
	{types.ApiDiagEndpointKey, "/slurm/v0.0.40/diag"},
	{types.ApiSharesEndpointKey, "/slurm/v0.0.40/shares"},
}
