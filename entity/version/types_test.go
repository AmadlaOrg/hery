package version

import (
	"regexp"
	"testing"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		uri           string
		expectedMatch bool
	}{
		{"github.com/user/repo@v1.0.0", true},
		{"github.com/user/repo@v2.3.4-beta.2", true},
		{"github.com/user/repo@v0.0.0-20170915032832-14c0d48ead0c", true},
		{"github.com/user/repo@v1.2.3", true},
		{"github.com/user/repo", false},
		{"github.com/user/repo@", false},
		{"github.com/user/repo@v1.0", false},
		{"github.com/user/repo@v1.0.0.0", false},
	}

	re := regexp.MustCompile(Match)
	for _, test := range tests {
		t.Run(test.uri, func(t *testing.T) {
			match := re.MatchString(test.uri)
			if match != test.expectedMatch {
				t.Errorf("Match(%s) = %v, want %v", test.uri, match, test.expectedMatch)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		version       string
		expectedMatch bool
	}{
		{"v1.0.0", true},
		{"v2.3.4-beta.2", true},
		{"v0.0.0-20170915032832-14c0d48ead0c", false}, // Not a valid semantic version
		{"v1.2.3", true},
		{"v1", true},
		{"v1.0", true},
		{"v1.0.0-rc.1", true},
		{"v1.0.0.0", false},
	}

	re := regexp.MustCompile(Format)
	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			match := re.MatchString(test.version)
			if match != test.expectedMatch {
				t.Errorf("Format(%s) = %v, want %v", test.version, match, test.expectedMatch)
			}
		})
	}
}

func TestFormatForDir(t *testing.T) {
	tests := []struct {
		uri             string
		expectedModule  string
		expectedVersion string
		expectedMatch   bool
	}{
		{"github.com/user/repo@v1.0.0", "github.com/user/repo", "v1.0.0", true},
		{"github.com/user/repo@v2.3.4-beta.2", "github.com/user/repo", "v2.3.4-beta.2", true},
		{"github.com/user/repo@v0.0.0-20170915032832-14c0d48ead0c", "github.com/user/repo", "v0.0.0-20170915032832-14c0d48ead0c", true},
		{"github.com/user/repo@v1.2.3", "github.com/user/repo", "v1.2.3", true},
		{"github.com/user/repo", "", "", false},
		{"github.com/user/repo@", "", "", false},
		{"github.com/user/repo@v1.0.0.0", "", "", false},
	}

	re := regexp.MustCompile(FormatForDir)
	for _, test := range tests {
		t.Run(test.uri, func(t *testing.T) {
			matches := re.FindStringSubmatch(test.uri)
			if (len(matches) == 0) != !test.expectedMatch {
				t.Errorf("FormatForDir(%s) = %v, want %v", test.uri, len(matches) > 0, test.expectedMatch)
			}
			if test.expectedMatch {
				if matches[1] != test.expectedModule {
					t.Errorf("FormatForDir(%s) module = %s, want %s", test.uri, matches[1], test.expectedModule)
				}
				if matches[2] != test.expectedVersion {
					t.Errorf("FormatForDir(%s) version = %s, want %s", test.uri, matches[2], test.expectedVersion)
				}
			}
		})
	}
}
