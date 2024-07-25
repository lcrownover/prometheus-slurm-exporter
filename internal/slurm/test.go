package slurm

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

func readTestDataBytes(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %v\n", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read file: %v\n", err)
	}

	return data
}

func unmarshalFixtureSlurmJobsResponse(d []byte) types.V0040OpenapiJobInfoResp {
	var r *types.V0040OpenapiJobInfoResp
	err := json.Unmarshal(d, r)
	if err != nil {
		log.Fatalf("failed to unmarshal json: %v\n", err)
	}
	return *r
}

