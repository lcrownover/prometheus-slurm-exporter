package slurm

import "testing"

func TestUnmarshalJobsResponse(t *testing.T) {
	fb := readTestDataBytes("testdata/V0040OpenapiJobInfoResp.json")
	_, err := UnmarshalJobsResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal jobs response: %v\n", err)
	}
}

func TestUnmarshalNodesResponse(t *testing.T) {
	fb := readTestDataBytes("testdata/V0040OpenapiNodesResp.json")
	_, err := UnmarshalNodesResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal nodes response: %v\n", err)
	}
}
