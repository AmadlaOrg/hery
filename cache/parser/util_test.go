package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJsonSchemaToSQLiteType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"string", "TEXT"},
		{"number", "BIGINT"},
		{"integer", "BIGINT"},
		{"object", "TEXT"},
		{"array", "TEXT"},
		{"boolean", "BOOLEAN"},
		{"null", "TEXT"},
		{"unknown", "TEXT"}, // Default case
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseJsonSchemaToSQLiteType(tt.input)
			if result != tt.expected {
				t.Errorf("parseJsonSchemaToSQLiteType(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseJsonSchemaFormatToSQLiteType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"date-time", "DATETIME"},
		{"time", "TEXT"},
		{"date", "DATE"},
		{"duration", "TEXT"},
		{"unknown", ""}, // Default case: returns an empty string
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseJsonSchemaFormatToSQLiteType(tt.input)
			if result != tt.expected {
				t.Errorf("parseJsonSchemaFormatToSQLiteType(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseJsonSchemaToSQLiteConstraint(t *testing.T) {
	tests := []struct {
		columnName           string
		jsonSchemaConstraint string
		value                any
		expected             string
	}{
		{"name", "minLength", 3, "LENGTH(name) >= 3"},
		{"description", "maxLength", 255, "LENGTH(description) <= 255"},
		{"status", "enum", []any{"active", "inactive"}, "CHECK (status IN ('active', 'inactive'))"},
		{"status", "enum", []any{}, ""}, // Empty enum should return empty constraint
		{"name", "unknown", 123, ""},    // Unknown constraint
	}

	for _, test := range tests {
		actual := parseJsonSchemaToSQLiteConstraint(test.columnName, test.jsonSchemaConstraint, test.value)
		assert.Equal(t, test.expected, actual)
	}
}

func TestJoinSQLiteConstraints(t *testing.T) {
	tests := []struct {
		constraints []string
		expected    string
	}{
		{[]string{"LENGTH(name) >= 3", "LENGTH(name) <= 255"}, "CHECK (LENGTH(name) >= 3 AND LENGTH(name) <= 255)"},
		{[]string{}, ""}, // No constraints
		{[]string{"LENGTH(name) >= 3"}, "CHECK (LENGTH(name) >= 3)"}, // Single constraint
	}

	for _, test := range tests {
		actual := joinSQLiteConstraints(test.constraints...)
		assert.Equal(t, test.expected, actual)
	}
}
