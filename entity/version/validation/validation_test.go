package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersionValidation_Exists(t *testing.T) {
	vs := SValidation{}

	tests := []struct {
		name     string
		version  string
		versions []string
		expected bool
	}{
		{
			name:     "Version exists in list",
			version:  "v1.0.0",
			versions: []string{"v1.0.0", "v2.0.0"},
			expected: true,
		},
		{
			name:     "Version does not exist in list",
			version:  "v1.1.0",
			versions: []string{"v1.0.0", "v2.0.0"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.Exists(tt.version, tt.versions)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersionValidation_Format(t *testing.T) {
	vs := SValidation{}

	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "Valid version v1.0.0",
			version:  "v1.0.0",
			expected: true,
		},
		{
			name:     "Valid version v1.0",
			version:  "v1.0",
			expected: false,
		},
		{
			name:     "Valid version v1",
			version:  "v1",
			expected: false,
		},
		{
			name:     "Invalid version format",
			version:  "1.0.0",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.Format(tt.version)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersionValidation_PseudoFormat(t *testing.T) {
	vs := SValidation{}

	tests := []struct {
		name          string
		pseudoVersion string
		expected      bool
	}{
		{
			name:          "Valid pseudo version",
			pseudoVersion: "v0.0.0-20230101120000-abcdef123456",
			expected:      true,
		},
		{
			name:          "Invalid pseudo version format",
			pseudoVersion: "v0.0.1-20230101120000-abcdef123456",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.PseudoFormat(tt.pseudoVersion)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersionValidation_AnyFormat(t *testing.T) {
	vs := SValidation{}

	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "Valid standard version v1.0.0",
			version:  "v1.0.0",
			expected: true,
		},
		{
			name:     "Valid pseudo version",
			version:  "v0.0.0-20230101120000-abcdef123456",
			expected: true,
		},
		{
			name:     "Valid latest version",
			version:  "latest",
			expected: true,
		},
		{
			name:     "Invalid version format",
			version:  "invalid-version",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vs.AnyFormat(tt.version)
			assert.Equal(t, tt.expected, result)
		})
	}
}
