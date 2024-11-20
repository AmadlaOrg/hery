package parser

import "testing"

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
