package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessRow(t *testing.T) {
	tests := []struct {
		name                 string
		inputRow             map[string]any
		expectedColumnNames  []string
		expectedPlaceholders []string
		expectedColumnValues []string
	}{
		{
			name:                 "empty row",
			inputRow:             map[string]any{},
			expectedColumnNames:  []string{},
			expectedPlaceholders: []string{},
			expectedColumnValues: []string{},
		},
		{
			name: "single column",
			inputRow: map[string]any{
				"Id": "123",
			},
			expectedColumnNames:  []string{"Id"},
			expectedPlaceholders: []string{"?"},
			expectedColumnValues: []string{"123"},
		},
		{
			name: "multiple columns",
			inputRow: map[string]any{
				"Id":    "123",
				"Name":  "John",
				"Age":   30,
				"Email": "john@example.com",
			},
			expectedColumnNames:  []string{"Id", "Name", "Age", "Email"},
			expectedPlaceholders: []string{"?", "?", "?", "?"},
			expectedColumnValues: []string{"123", "John", "30", "john@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columnNames, placeholders, columnValues := processRow(tt.inputRow)

			assert.ElementsMatch(t, tt.expectedColumnNames, columnNames, "Column names mismatch")
			assert.ElementsMatch(t, tt.expectedPlaceholders, placeholders, "Placeholders mismatch")
			assert.ElementsMatch(t, tt.expectedColumnValues, columnValues, "Column values mismatch")
		})
	}
}

func TestProcessWhere(t *testing.T) {
	tests := []struct {
		name       string
		inputWhere map[string]any
		expected   []string
	}{
		{
			name:       "empty",
			inputWhere: make(map[string]any),
			expected:   []string{},
		},
		{
			name: "single",
			inputWhere: map[string]any{
				"Id":          "c6beaec1-90c4-4d2a-aaef-211ab00b86bd",
				"server_name": "localhost",
				"listen":      "[80, 443]",
			},
			expected: []string{
				"Id = 'c6beaec1-90c4-4d2a-aaef-211ab00b86bd'",
				"server_name = 'localhost'",
				"listen = '[80, 443]'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processWhere(tt.inputWhere)
			assert.ElementsMatch(t, tt.expected, got)
		})
	}
}
