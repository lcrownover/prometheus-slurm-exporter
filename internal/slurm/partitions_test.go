package slurm
//
// import (
// 	"testing"
//
// 	"github.com/lcrownover/prometheus-slurm-exporter/internal/api"
// 	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
// )
//
// func TestParsePartitionsMetrics(t *testing.T) {
// 	partitionsBytes := util.ReadTestDataBytes("V0040OpenapiPartitionResp.json")
// 	partitionResp, _ := api.UnmarshalPartitionsResponse(partitionsBytes)
// 	_, err := ParsePartitionsMetrics(*partitionResp)
// 	if err != nil {
// 		t.Fatalf("failed to parse partitions metrics: %v\n", err)
// 	}
// }
