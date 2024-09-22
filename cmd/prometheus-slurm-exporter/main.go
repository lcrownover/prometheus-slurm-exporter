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

	apiCacheTimeoutSecondsStr, found := os.LookupEnv("SLURM_EXPORTER_API_CACHE_TIMEOUT_SECONDS")
	if !found {
		apiCacheTimeoutSecondsStr = "5"
	}
	apiCacheTimeout, err := api.ParseCacheTimeoutSeconds(apiCacheTimeoutSecondsStr)
	if err != nil {
		fmt.Println("Invalid value for SLURM_EXPORTER_API_CACHE_TIMEOUT_SECONDS")
		os.Exit(1)
	}

	// Create the API Cache
	cache := api.NewApiCache(apiCacheTimeout)

	// Set up the context to pass around
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.ApiUserKey, apiUser)
	ctx = context.WithValue(ctx, types.ApiTokenKey, apiToken)
	ctx = context.WithValue(ctx, types.ApiURLKey, apiURL)
	ctx = context.WithValue(ctx, types.ApiCacheKey, &cache)

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

	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	log.Printf("Starting Server: %s\n", listenAddress)
	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
