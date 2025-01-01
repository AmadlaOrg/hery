package pointer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToPtr(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "String to Pointer",
			input:    "test string",
			expected: "test string",
		},
		{
			name:     "Integer to Pointer",
			input:    42,
			expected: 42,
		},
		{
			name:     "Boolean to Pointer",
			input:    true,
			expected: true,
		},
		{
			name:     "Float to Pointer",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "Struct to Pointer",
			input:    struct{ Field string }{Field: "value"},
			expected: struct{ Field string }{Field: "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToPtr(tt.input)
			if got == nil {
				t.Fatalf("ToPtr(%v) returned nil, expected a pointer", tt.input)
			}
			assert.Equal(t, tt.expected, *got)
		})
	}
}
