package api

import (
	"testing"
)

func TestCleanseBaseURL(t *testing.T) {
	tts := []struct {
		in   string
		want string
	}{
		{"https://google.com", "google.com"},
		{"http://google.com", "google.com"},
		{"google.com", "google.com"},
	}
	for _, tt := range tts {
		t.Run(tt.in, func(t *testing.T) {
			got := CleanseBaseURL(tt.in)
			if got != tt.want {
				t.Fatalf("failed to cleanse base url: got %q, want %q", got, tt.want)
			}
		})
	}
}
