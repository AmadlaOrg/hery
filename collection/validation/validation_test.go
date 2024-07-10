package validation

import "testing"

func TestName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		hasError bool
	}{
		{"ValidName", "valid_name-123", true, false},
		{"ValidNameWithUnderscore", "valid_name", true, false},
		{"ValidNameWithDash", "valid-name", true, false},
		{"ValidNameWithNumbers", "valid123", true, false},
		{"InvalidNameWithSpace", "invalid name", false, false},
		{"InvalidNameWithSpecialChar", "invalid@name", false, false},
		{"EmptyName", "", false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Name(test.input)
			if result != test.expected {
				t.Errorf("CollectionName(%q) = %v; expected %v", test.input, result, test.expected)
			}
		})
	}
}
