package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateColumnName(t *testing.T) {
	tests := []struct {
		name        string
		columnName  string
		expectedErr string
	}{
		// Valid cases
		{name: "Valid column name", columnName: "username", expectedErr: ""},
		{name: "Valid column name with underscore", columnName: "user_name", expectedErr: ""},
		{name: "Valid column name starting with underscore", columnName: "_username", expectedErr: ""},

		// Invalid cases
		{name: "Empty column name", columnName: "", expectedErr: "column name cannot be empty"},
		{name: "Column name with spaces", columnName: "user name", expectedErr: "column name 'user name' contains invalid characters"},
		{name: "Column name starting with number", columnName: "1username", expectedErr: "column name '1username' contains invalid characters"},
		{name: "Column name with special characters", columnName: "user@name", expectedErr: "column name 'user@name' contains invalid characters"},
		{name: "Column name is reserved keyword", columnName: "SELECT", expectedErr: "column name 'SELECT' is a reserved keyword"},
		{name: "Column name is reserved keyword (case insensitive)", columnName: "select", expectedErr: "column name 'select' is a reserved keyword"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColumnName(tt.columnName)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}

func TestValidateOperator(t *testing.T) {
	tests := []struct {
		name      string
		operator  string
		expectErr bool
	}{
		{name: "Valid operator - equals", operator: "=", expectErr: false},
		{name: "Valid operator - not equals", operator: "<>", expectErr: false},
		{name: "Valid operator - less than", operator: "<", expectErr: false},
		{name: "Valid operator - greater than", operator: ">", expectErr: false},
		{name: "Valid operator - LIKE", operator: "LIKE", expectErr: false},
		{name: "Invalid operator", operator: "INVALID", expectErr: true},
		{name: "Empty operator", operator: "", expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOperator(tt.operator)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid operator")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

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
			got := buildWhere(tt.inputWhere)
			assert.ElementsMatch(t, tt.expected, got)
		})
	}
}

func TestBuildGroupBy(t *testing.T) {
	tests := []struct {
		name     string
		groupBy  map[string]any
		expected string
	}{
		{name: "Empty GroupBy", groupBy: map[string]any{}, expected: ""},
		{name: "Single GroupBy", groupBy: map[string]any{"column1": "ASC"}, expected: " GROUP BY column1 ASC"},
		{name: "Multiple GroupBy", groupBy: map[string]any{"column1": "ASC", "column2": "DESC"}, expected: " GROUP BY column1 ASC column2 DESC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildGroupBy(tt.groupBy)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildHaving(t *testing.T) {
	tests := []struct {
		name     string
		having   map[string]any
		expected string
	}{
		{name: "Empty Having", having: map[string]any{}, expected: ""},
		{name: "Single Having", having: map[string]any{"column1": "value1"}, expected: " HAVING column1 = 'value1'"},
		{name: "Multiple Having", having: map[string]any{"column1": "value1", "column2": "value2"}, expected: " HAVING column1 = 'value1' column2 = 'value2'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildHaving(tt.having)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildOrderBy(t *testing.T) {
	tests := []struct {
		name     string
		orderBy  map[string]any
		expected string
	}{
		{name: "Empty OrderBy", orderBy: map[string]any{}, expected: ""},
		{name: "Single OrderBy", orderBy: map[string]any{"column1": "ASC"}, expected: " ORDER BY column1 ASC"},
		{name: "Multiple OrderBy", orderBy: map[string]any{"column1": "ASC", "column2": "DESC"}, expected: " ORDER BY column1 ASC column2 DESC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildOrderBy(tt.orderBy)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildLimit(t *testing.T) {
	tests := []struct {
		name     string
		limit    int64
		expected string
	}{
		{name: "Valid limit", limit: 10, expected: " LIMIT 10"},
		{name: "Zero limit", limit: 0, expected: ""},
		{name: "Negative limit", limit: -5, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildLimit(tt.limit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildOffset(t *testing.T) {
	tests := []struct {
		name     string
		offset   int64
		expected string
	}{
		{name: "Valid offset", offset: 20, expected: " OFFSET 20"},
		{name: "Zero offset", offset: 0, expected: ""},
		{name: "Negative offset", offset: -10, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildOffset(tt.offset)
			assert.Equal(t, tt.expected, result)
		})
	}
}
