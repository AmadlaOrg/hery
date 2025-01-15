package parser

import (
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/entity/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJsonSchemaToSQLiteType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"string", "TEXT"},
		{"number", "NUMERIC"},
		{"integer", "NUMERIC"},
		{"object", "TEXT"},
		{"array", "TEXT"},
		{"boolean", "BOOLEAN"},
		{"null", "TEXT"},
		{"unknown", "TEXT"}, // Default case
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseJsonSchemaToSQLiteType(schema.DataType(tt.input))
			assert.Equal(t, database.DataType(tt.expected), result)
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
		// TODO: Is this good idea to defautl to TEXT?
		{"unknown", "TEXT"}, // Default case: returns an empty string
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseJsonSchemaFormatToSQLiteType(schema.DataFormat(tt.input))
			assert.Equal(t, database.DataType(tt.expected), result)
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
