package validation

import (
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		version  string
		expected string
		isValid  bool
	}{
		{"v1.0.0", "v1.0.0", true},
		{"v1.0", "v1.0", true},
		{"v1", "v1", true},
		{"v1.0.0-beta", "", false},
		{"1.0.0", "", false},
		{"v1.0.0.0", "", false},
	}

	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			validated, err := Format(test.version)
			if test.isValid {
				if err != nil {
					t.Errorf("Format(%s) returned an error: %v", test.version, err)
				} else if validated != test.expected {
					t.Errorf("Format(%s) = %s; expected %s", test.version, validated, test.expected)
				}
			} else {
				if err == nil {
					t.Errorf("Format(%s) expected an error, but got none", test.version)
				}
			}
		})
	}
}
