package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityPathUrl(t *testing.T) {
	tests := []struct {
		entityUrl string
		version   string
		expected  string
	}{
		{"https://github.com/user/repo@v1.0.0", "v1.0.0", "https://github.com/user/repo"},
		{"https://github.com/user/repo@latest", "latest", "https://github.com/user/repo"},
		{"https://github.com/user/repo@v2.3.4", "v2.3.4", "https://github.com/user/repo"},
		{"https://github.com/user/repo", "v1.0.0", "https://github.com/user/repo"},
	}

	for _, test := range tests {
		t.Run(test.entityUrl, func(t *testing.T) {
			result := EntityPathUrl(test.entityUrl, test.version)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestEntityFullRepoUrl(t *testing.T) {
	tests := []struct {
		entityUrl string
		expected  string
	}{
		{"github.com/user/repo", "https://github.com/user/repo"},
		{"gitlab.com/user/repo", "https://gitlab.com/user/repo"},
		{"bitbucket.org/user/repo", "https://bitbucket.org/user/repo"},
		{"example.com/user/repo", "https://example.com/user/repo"},
	}

	for _, test := range tests {
		t.Run(test.entityUrl, func(t *testing.T) {
			result := EntityFullRepoUrl(test.entityUrl)
			assert.Equal(t, test.expected, result)
		})
	}
}
