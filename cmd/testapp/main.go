package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

func main() {
	fmt.Println("trying api json")
	b, err := os.ReadFile("/tmp/shares-api.json")
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
	}
	var ra types.V0040OpenapiSharesResp
	err = json.Unmarshal(b, &ra)
	if err != nil {
		fmt.Printf("failed to unmarshal json: %v\n", err)
	}
	fmt.Printf("%+v\n", ra)

	fmt.Println("trying copied json")
	b, err = os.ReadFile("/tmp/shares-copy.json")
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
	}
	var rc types.V0040OpenapiSharesResp
	err = json.Unmarshal(b, &rc)
	if err != nil {
		fmt.Printf("failed to unmarshal json: %v\n", err)
	}
	fmt.Printf("%+v\n", rc)
}
