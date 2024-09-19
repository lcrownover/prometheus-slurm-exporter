package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/slurm"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var err error

func main() {

	// set up logging
	lvl := slog.LevelInfo
	_, found := os.LookupEnv("SLURM_EXPORTER_DEBUG")
	if found {
		lvl = slog.LevelDebug
	}
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(l)
	slog.Debug("debug logging enabled")

	listenAddress, found := os.LookupEnv("SLURM_EXPORTER_LISTEN_ADDRESS")
	if !found {
		listenAddress = ":8080"
	}

	apiUser, found := os.LookupEnv("SLURM_EXPORTER_API_USER")
	if !found {
		fmt.Println("You must set SLURM_EXPORTER_API_USER")
		os.Exit(1)
	}

	apiToken, found := os.LookupEnv("SLURM_EXPORTER_API_TOKEN")
	if !found {
		fmt.Println("You must set SLURM_EXPORTER_API_TOKEN")
		os.Exit(1)
	}

	apiURL, found := os.LookupEnv("SLURM_EXPORTER_API_URL")
	if !found {
		fmt.Println("You must set SLURM_EXPORTER_API_URL. Example: localhost:6820")
		os.Exit(1)
	}
	apiURL = api.CleanseBaseURL(apiURL)

	// Set up the context to pass around
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.ApiUserKey, apiUser)
	ctx = context.WithValue(ctx, types.ApiTokenKey, apiToken)
	ctx = context.WithValue(ctx, types.ApiURLKey, apiURL)

	// Register all the endpoints
	ctx = context.WithValue(ctx, types.ApiJobsEndpointKey, "/slurm/v0.0.40/jobs")
	ctx = context.WithValue(ctx, types.ApiNodesEndpointKey, "/slurm/v0.0.40/nodes")
	ctx = context.WithValue(ctx, types.ApiPartitionsEndpointKey, "/slurm/v0.0.40/partitions")
	ctx = context.WithValue(ctx, types.ApiDiagEndpointKey, "/slurm/v0.0.40/diag")
	ctx = context.WithValue(ctx, types.ApiSharesEndpointKey, "/slurm/v0.0.40/shares")

	r := prometheus.NewRegistry()

	// r.MustRegister(slurm.NewAccountsCollector(ctx))
	// r.MustRegister(slurm.NewCPUsCollector(ctx))
	// r.MustRegister(slurm.NewGPUsCollector(ctx))
	// r.MustRegister(slurm.NewNodesCollector(ctx))
	// r.MustRegister(slurm.NewNodeCollector(ctx))

	// TODO: write and test this
	// r.MustRegister(slurm.NewPartitionsCollector(ctx))
	// r.MustRegister(slurm.NewPartitionsCollectorOld())

	// TODO: json parsing bug
	r.MustRegister(slurm.NewFairShareCollector(ctx))
	// r.MustRegister(slurm.NewFairShareCollectorOld())

	// r.MustRegister(slurm.NewQueueCollector(ctx))
	// r.MustRegister(slurm.NewSchedulerCollector(ctx))
	// r.MustRegister(slurm.NewUsersCollector(ctx))

	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	log.Printf("Starting Server: %s\n", listenAddress)
	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
