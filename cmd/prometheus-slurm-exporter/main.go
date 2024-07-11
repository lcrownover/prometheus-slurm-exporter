/* Copyright 2017-2022 Lucas Crownover, Victor Penso, Matteo Dessalvi, Joeri Hermans

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/slurm"
	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
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
	apiURL = util.CleanseBaseURL(apiURL)

	ctx := context.Background()
	ctx = context.WithValue(ctx, util.ApiUserKey, apiUser)
	ctx = context.WithValue(ctx, util.ApiTokenKey, apiToken)
	ctx = context.WithValue(ctx, util.ApiURLKey, apiURL)

	prometheus.MustRegister(slurm.NewAccountsCollector(ctx)) // from accounts.go
	prometheus.MustRegister(slurm.NewOldAccountsCollector()) // from accounts.go
	// prometheus.MustRegister(slurm.NewCPUsCollector())       // from cpus.go
	// prometheus.MustRegister(slurm.NewNodesCollector())      // from nodes.go
	// prometheus.MustRegister(slurm.NewNodeCollector())       // from node.go
	// prometheus.MustRegister(slurm.NewPartitionsCollector()) // from partitions.go
	// prometheus.MustRegister(slurm.NewQueueCollector())      // from queue.go
	// prometheus.MustRegister(slurm.NewSchedulerCollector())  // from scheduler.go
	// prometheus.MustRegister(slurm.NewFairShareCollector())  // from sshare.go
	// prometheus.MustRegister(slurm.NewUsersCollector())      // from users.go

	// gpuAcctString := os.Getenv("SLURM_EXPORTER_GPU_ACCOUNTING")
	// if gpuAcctString == "true" || gpuAcctString == "1" {
	// 	prometheus.MustRegister(slurm.NewGPUsCollector())
	// 	log.Println("GPUs Accounting ON")
	// }

	log.Printf("Starting Server: %s\n", listenAddress)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
