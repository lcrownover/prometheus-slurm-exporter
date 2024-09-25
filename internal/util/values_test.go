package util

import (
	"testing"
)

func TestGetValueOrZero(t *testing.T) {
	tests := []struct {
		name     string
		input    any // using any to support both *int32 and *int64
		expected float64
	}{
		{
			name:     "nil int32 pointer",
			input:    (*int32)(nil),
			expected: 0,
		},
		{
			name:     "non-nil int32 pointer",
			input:    func() *int32 { v := int32(32); return &v }(),
			expected: 32,
		},
		{
			name:     "nil int64 pointer",
			input:    (*int64)(nil),
			expected: 0,
		},
		{
			name:     "non-nil int64 pointer",
			input:    func() *int64 { v := int64(64); return &v }(),
			expected: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result float64
			switch v := tt.input.(type) {
			case *int32:
				result = GetValueOrZero(v)
			case *int64:
				result = GetValueOrZero(v)
			default:
				t.Fatalf("unsupported type: %T", v)
			}

			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
