package slurm

import (
	"testing"

	"github.com/lcrownover/prometheus-slurm-exporter/internal/util"
)

func TestCleanseBaseURL(t *testing.T) {
	tts := []struct{
		in string
		want string
	}{
		{"https://google.com", "google.com"},
		{"http://google.com", "google.com"},
		{"google.com", "google.com"},
	}
	for _,tt := range tts {
		t.Run(tt.in, func (t *testing.T) {
			got := CleanseBaseURL(tt.in)
			if got != tt.want {
				t.Fatalf("failed to cleanse base url: got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnmarshalJobsResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiJobInfoResp.json")
	_, err := UnmarshalJobsResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal jobs response: %v\n", err)
	}
}

func TestUnmarshalNodesResponse(t *testing.T) {
	fb := util.ReadTestDataBytes("V0040OpenapiNodesResp.json")
	_, err := UnmarshalNodesResponse(fb)
	if err != nil {
		t.Fatalf("failed to unmarshal nodes response: %v\n", err)
	}
}
