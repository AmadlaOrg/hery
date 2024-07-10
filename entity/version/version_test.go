package version

import "testing"

func TestExtract(t *testing.T) {
	tests := []struct {
		url           string
		expected      string
		expectedError bool
	}{
		{"github.com/user/repo@v1.0.0", "v1.0.0", false},
		{"github.com/user/repo@v2.3.4", "v2.3.4", false},
		{"github.com/user/repo", "", true},
		{"github.com/user/repo@", "", true},
		{"github.com/user/repo@v1.2.3-beta", "v1.2.3-beta", false},
		{"github.com/user/repo@v1.2.3-beta.2", "v1.2.3-beta.2", false},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			version, err := Extract(test.url)
			if (err != nil) != test.expectedError {
				t.Errorf("Extract(%s) returned error: %v, expectedError: %v", test.url, err, test.expectedError)
			}
			if version != test.expected {
				t.Errorf("Extract(%s) = %s, want %s", test.url, version, test.expected)
			}
		})
	}
}
