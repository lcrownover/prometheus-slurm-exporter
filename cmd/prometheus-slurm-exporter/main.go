package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akyoto/cache"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/slurm"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
	"github.com/prometheus/client_golang/prometheus"
)

var err error

var version = "2.1.1-beta"

func main() {
	// set up logging
	lvl := slog.LevelInfo
	_, debug := os.LookupEnv("SLURM_EXPORTER_DEBUG")
	if debug {
		lvl = slog.LevelDebug
	}
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(l)
	slog.Debug("debug logging enabled")

	// if -v is passed, print the version and exit
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Printf("Starting Prometheus Slurm Exporter %s\n", version)

	listenAddress, found := os.LookupEnv("SLURM_EXPORTER_LISTEN_ADDRESS")
	if !found {
		listenAddress = "0.0.0.0:8080"
	}

	apiURL, found := os.LookupEnv("SLURM_EXPORTER_API_URL")
	if !found {
		fmt.Println("You must set SLURM_EXPORTER_API_URL. Example: localhost:6820")
		os.Exit(1)
	}

	var apiUser string
	var apiToken string
	var tlsEnable bool
	var tlsCert string
	var tlsKey string

	// we only need these values if the endpoint is not unix://
	if strings.HasPrefix(apiURL, "http://") || strings.HasPrefix(apiURL, "https://") {
		var found bool
		apiUser, found = os.LookupEnv("SLURM_EXPORTER_API_USER")
		if !found {
			fmt.Println("You must set SLURM_EXPORTER_API_USER")
			os.Exit(1)
		}

		apiToken, found = os.LookupEnv("SLURM_EXPORTER_API_TOKEN")
		if !found {
			fmt.Println("You must set SLURM_EXPORTER_API_TOKEN")
			os.Exit(1)
		}

		tlsString, found := os.LookupEnv("SLURM_EXPORTER_ENABLE_TLS")

		if !found {
			tlsEnable = false // default to false, do not break existing conf files
		} else {
			tlsEnable, err = strconv.ParseBool(tlsString)
			if err != nil {
				fmt.Println("Failed to parse SLURM_EXPORTER_ENABLE_TLS.  Please set to 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, or False.")
			}
		}
		if tlsEnable { // require tlsCert and tlsKey only if tlsEnable is true
			tlsCert, found = os.LookupEnv("SLURM_EXPORTER_TLS_CERT_PATH")
			if !found {
				fmt.Println("You must set SLURM_EXPORTER_TLS_CERT_PATH to the path of your cert")
				os.Exit(1)
			}
			tlsKey, found = os.LookupEnv("SLURM_EXPORTER_TLS_KEY_PATH")
			if !found {
				fmt.Println("You must set SLURM_EXPORTER_TLS_KEY_PATH to the path of your key")
				os.Exit(1)
			}
		}

	} else if strings.HasPrefix(apiURL, "unix://") {
		apiUser = ""
		apiToken = ""
		tlsEnable = false
		tlsCert = ""
		tlsKey = ""

	} else {
		fmt.Println("SLURM_EXPORTER_API_URL must start with unix://, http://, or https://")
		fmt.Println("Got: ", apiURL)
		os.Exit(1)
	}
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
	if tlsEnable {
		log.Fatal(http.ListenAndServeTLS(listenAddress, tlsCert, tlsKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(listenAddress, nil))
	}
}
