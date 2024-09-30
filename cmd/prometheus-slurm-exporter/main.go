package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/slurm"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

var err error

var version = "1.0.9"

func main() {
	// if -v is passed, print the version and exit
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Printf("Starting Prometheus Slurm Exporter %s\n", version)

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
		listenAddress = "0.0.0.0:8080"
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

	// API Cache
	apiCache := cache.New(60 * time.Second)

	// Set up the context to pass around
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.ApiUserKey, apiUser)
	ctx = context.WithValue(ctx, types.ApiTokenKey, apiToken)
	ctx = context.WithValue(ctx, types.ApiURLKey, apiURL)
	ctx = context.WithValue(ctx, types.ApiCacheKey, apiCache)

	// Register all the endpoints
	ctx = api.RegisterEndpoints(ctx)

	// Register all the collectors
	r := prometheus.NewRegistry()
	r.MustRegister(slurm.NewAccountsCollector(ctx))
	r.MustRegister(slurm.NewCPUsCollector(ctx))
	r.MustRegister(slurm.NewGPUsCollector(ctx))
	r.MustRegister(slurm.NewNodesCollector(ctx))
	r.MustRegister(slurm.NewNodeCollector(ctx))
	r.MustRegister(slurm.NewPartitionsCollector(ctx))
	r.MustRegister(slurm.NewFairShareCollector(ctx))
	r.MustRegister(slurm.NewQueueCollector(ctx))
	r.MustRegister(slurm.NewSchedulerCollector(ctx))
	r.MustRegister(slurm.NewUsersCollector(ctx))

	log.Printf("Starting Server: %s\n", listenAddress)
	http.Handle("/metrics", api.MetricsHandler(r, ctx))
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
