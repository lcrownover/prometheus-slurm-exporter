package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

func main() {
	b, err := os.ReadFile("/tmp/shares.json")
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
	}
	var r types.V0040OpenapiSharesResp
	err = json.Unmarshal(b, &r)
	if err != nil {
		fmt.Println("failed to unmarshal json: %v\n", err)
	}
	fmt.Printf("%+v\n", r)
}
