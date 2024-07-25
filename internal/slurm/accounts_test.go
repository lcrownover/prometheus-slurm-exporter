package slurm

import "testing"

func TestParseAccountMetrics(t *testing.T) {
	fb := readTestDataBytes("testdata/V0040OpenapiJobInfoResp.json")
	jobsResp, _ := UnmarshalJobsResponse(fb)
	m, err := ParseAccountsMetrics(jobsResp.Jobs)	
	if err != nil {
		t.Fatalf("failed to parse account metrics: %v\n", err)
	}
	m
}
