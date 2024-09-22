package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func beforeCollect(ctx context.Context) context.Context {
	slog.Info("Before collecting metrics")
	c := cache.New(60 * time.Second)
	ctx = context.WithValue(ctx, types.ApiCacheKey, c)
	err := PopulateCache(ctx)
	if err != nil {
		slog.Error("error populating request cache", "error", err)
	}
	return ctx
}

func afterCollect(ctx context.Context) {
	slog.Info("After collecting metrics")
	cache := ctx.Value(types.ApiCacheKey).(*cache.Cache)
	cache.Close()
}

func MetricsHandler(r *prometheus.Registry, ctx context.Context) http.HandlerFunc {
	h := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	return func(w http.ResponseWriter, r *http.Request) {
		ctx = beforeCollect(ctx)
		c := ctx.Value(types.ApiCacheKey).(*cache.Cache)
		v, _ := c.Get("diag")
		slog.Info("diag out", "v", v)
		h.ServeHTTP(w, r)
		afterCollect(ctx)
	}
}
