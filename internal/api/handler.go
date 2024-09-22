package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func beforeCollect(ctx context.Context) {
	err := PopulateCache(ctx)
	if err != nil {
		slog.Error("error populating request cache", "error", err)
	}
}

func afterCollect(ctx context.Context) {
	WipeCache(ctx)
}

func MetricsHandler(r *prometheus.Registry, ctx context.Context) http.HandlerFunc {
	h := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	return func(w http.ResponseWriter, r *http.Request) {
		beforeCollect(ctx)
		h.ServeHTTP(w, r)
		afterCollect(ctx)
	}
}
