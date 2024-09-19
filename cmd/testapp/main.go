package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/types"
)

func main() {
	b, _ := os.ReadFile("/tmp/shares.json")
	var r types.V0040OpenapiSharesResp
	err := json.Unmarshal(b, &r)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", r)
}
