//go:build 2311

package api

import (
	"context"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

func RegisterEndpoints(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, types.ApiJobsEndpointKey, "/slurm/v0.0.40/jobs")
	ctx = context.WithValue(ctx, types.ApiNodesEndpointKey, "/slurm/v0.0.40/nodes")
	ctx = context.WithValue(ctx, types.ApiPartitionsEndpointKey, "/slurm/v0.0.40/partitions")
	ctx = context.WithValue(ctx, types.ApiDiagEndpointKey, "/slurm/v0.0.40/diag")
	ctx = context.WithValue(ctx, types.ApiSharesEndpointKey, "/slurm/v0.0.40/shares")
	return ctx
}
