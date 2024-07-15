package version

import (
	"testing"
)

func TestExtract(t *testing.T) {
	/*tests := []struct {
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
	}*/
}

func TestExists(t *testing.T) {
	// Mock the remote.Tags function
	/*remote.Tags = func(entityUrlPath string) ([]string, error) {
		return []string{"v1.0.0", "v2.0.0", "v1.1.0"}, nil
	}

	tests := []struct {
		entityUrlPath string
		version       string
		expected      bool
		err           bool
	}{
		{"some/url", "v1.0.0", true, false},
		{"some/url", "v2.0.0", true, false},
		{"some/url", "v3.0.0", false, false},
		{"some/url", "v1.1.0", true, false},
	}

	for _, test := range tests {
		result, err := Exists(test.entityUrlPath, test.version)
		if (err != nil) != test.err {
			t.Errorf("Exists(%v, %v) unexpected error: %v", test.entityUrlPath, test.version, err)
		}
		if result != test.expected {
			t.Errorf("Exists(%v, %v) = %v; want %v", test.entityUrlPath, test.version, result, test.expected)
		}
	}*/
}

func TestLatest(t *testing.T) {
	/*tests := []struct {
		versions []string
		expected string
		err      bool
	}{
		{[]string{"v1.0.0", "v2.0.0", "v1.1.0"}, "v2.0.0", false},
		{[]string{"v1.0.0", "v1.0.1", "v1.0.2"}, "v1.0.2", false},
		{[]string{"v1.0.0-alpha.1", "v1.0.0-beta.1", "v1.0.0-rc.1", "v1.0.0"}, "v1.0.0", false},
		{[]string{}, "", true},
	}

	for _, test := range tests {
		result, err := Latest(test.versions)
		if (err != nil) != test.err {
			t.Errorf("Latest(%v) unexpected error: %v", test.versions, err)
		}
		if result != test.expected {
			t.Errorf("Latest(%v) = %v; want %v", test.versions, result, test.expected)
		}
	}*/
}

func TestVersionLess(t *testing.T) {
	/*tests := []struct {
		v1, v2   string
		expected bool
	}{
		{"v1.0.0", "v2.0.0", true},
		{"v2.0.0", "v1.0.0", false},
		{"v1.0.0", "v1.0.0", false},
		{"v1.0.0-alpha.1", "v1.0.0-beta.1", true},
		{"v1.0.0-rc.1", "v1.0.0", true},
	}

	for _, test := range tests {
		result := versionLess(test.v1, test.v2)
		if result != test.expected {
			t.Errorf("VersionLess(%v, %v) = %v; want %v", test.v1, test.v2, result, test.expected)
		}
	}*/
}

func TestCompareVersions(t *testing.T) {
	/*tests := []struct {
		v1, v2   string
		expected int
	}{
		{"v1.0.0", "v2.0.0", -1},
		{"v2.0.0", "v1.0.0", 1},
		{"v1.0.0", "v1.0.0", 0},
		{"v1.0.0-alpha.1", "v1.0.0-beta.1", -1},
		{"v1.0.0-rc.1", "v1.0.0", -1},
	}

	for _, test := range tests {
		result := compareVersions(test.v1, test.v2)
		if result != test.expected {
			t.Errorf("CompareVersions(%v, %v) = %v; want %v", test.v1, test.v2, result, test.expected)
		}
	}*/
}

func TestParseVersion(t *testing.T) {
	/*tests := []struct {
		version  string
		expected []int
	}{
		{"v1.0.0", []int{1, 0, 0}},
		{"v2.1.3", []int{2, 1, 3}},
		{"v1.0.0-alpha.1", []int{1, 0, 0}},
	}

	for _, test := range tests {
		result, _ := parseVersion(test.version)
		for i := range test.expected {
			if result[i] != test.expected[i] {
				t.Errorf("ParseVersion(%v) = %v; want %v", test.version, result, test.expected)
			}
		}
	}*/
}
