package validation

import (
	"testing"
)

func TestEntityUri(t *testing.T) {
	entityValidationService := NewEntityValidationService()
	tests := []struct {
		path     string
		expected bool
	}{
		{"github.com/user/module", true},
		{"github.com/user/module@v1.0.0", true},
		{"https://github.com/user/module", false},
		{"github.com/user/module with spaces", false},
		{"github.com/user/module:colon", false},
		{"github.com/user/module?query", false},
		{"github.com/user/module&param", false},
		{"github.com/user/module=equal", false},
		{"github.com/user/module#fragment", false},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			result := entityValidationService.EntityUri(test.path)
			if result != test.expected {
				t.Errorf("EntityUrl(%q) = %v; expected %v", test.path, result, test.expected)
			}
		})
	}
}
