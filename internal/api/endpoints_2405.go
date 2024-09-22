//go:build 2405

package api

import (
	"context"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

func RegisterEndpoints(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, types.ApiJobsEndpointKey, "/slurm/v0.0.41/jobs")
	ctx = context.WithValue(ctx, types.ApiNodesEndpointKey, "/slurm/v0.0.41/nodes")
	ctx = context.WithValue(ctx, types.ApiPartitionsEndpointKey, "/slurm/v0.0.41/partitions")
	ctx = context.WithValue(ctx, types.ApiDiagEndpointKey, "/slurm/v0.0.41/diag")
	ctx = context.WithValue(ctx, types.ApiSharesEndpointKey, "/slurm/v0.0.41/shares")
	return ctx
}
