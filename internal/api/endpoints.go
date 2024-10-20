package api

import (
	"context"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

type endpoint struct {
	key  types.Key
	path string
}

// this gives a compile warning but centralizes the endpoints
var endpoints = versionedEndpoints

func RegisterEndpoints(ctx context.Context) context.Context {
	for _, e := range endpoints {
		ctx = context.WithValue(ctx, e.key, e.path)
	}
	return ctx
}
