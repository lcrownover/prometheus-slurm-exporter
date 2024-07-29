package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

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
	apiURL = slurm.CleanseBaseURL(apiURL)

	// Set up the context to pass around
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.ApiUserKey, apiUser)
	ctx = context.WithValue(ctx, types.ApiTokenKey, apiToken)
	ctx = context.WithValue(ctx, types.ApiURLKey, apiURL)

	// Register all the endpoints
	ctx = context.WithValue(ctx, types.ApiJobsEndpointKey, "/slurm/v0.0.40/jobs")
	ctx = context.WithValue(ctx, types.ApiNodesEndpointKey, "/slurm/v0.0.40/nodes")

	r := prometheus.NewRegistry()

	// r.MustRegister(slurm.NewAccountsCollector(ctx)) // from accounts.go
	// r.MustRegister(slurm.NewOldAccountsCollector()) // from accounts.go

	r.MustRegister(slurm.NewCPUsCollector(ctx)) // from cpus.go
	r.MustRegister(slurm.NewCPUsCollectorOld()) // from cpus.go

	// r.MustRegister(slurm.NewNodesCollector())      // from nodes.go
	// r.MustRegister(slurm.NewNodeCollector())       // from node.go
	// r.MustRegister(slurm.NewPartitionsCollector()) // from partitions.go
	// r.MustRegister(slurm.NewQueueCollector())      // from queue.go
	// r.MustRegister(slurm.NewSchedulerCollector())  // from scheduler.go
	// r.MustRegister(slurm.NewFairShareCollector())  // from sshare.go
	// r.MustRegister(slurm.NewUsersCollector())      // from users.go

	// gpuAcctString := os.Getenv("SLURM_EXPORTER_GPU_ACCOUNTING")
	// if gpuAcctString == "true" || gpuAcctString == "1" {
	// 	r.MustRegister(slurm.NewGPUsCollector())
	// 	log.Println("GPUs Accounting ON")
	// }
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	log.Printf("Starting Server: %s\n", listenAddress)
	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
