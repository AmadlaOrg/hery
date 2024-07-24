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

func TestExtractRepoPath(t *testing.T) {
	tests := []struct {
		repoURL       string
		expectedPath  string
		expectedError bool
	}{
		{
			repoURL:       "github.com/AmadlaOrg/EntityApplication/Network/Web/Server",
			expectedPath:  "https://github.com/AmadlaOrg/EntityApplication",
			expectedError: false,
		},
		{
			repoURL:       "github.com/AmadlaOrg/EntityApplication/WebServer",
			expectedPath:  "https://github.com/AmadlaOrg/EntityApplication",
			expectedError: false,
		},
		{
			repoURL:       "git.suckless.org/st/",
			expectedPath:  "https://git.suckless.org/st",
			expectedError: false,
		},
		{
			repoURL:       "github.com/owner/repo",
			expectedPath:  "https://github.com/owner/repo",
			expectedError: false,
		},
		{
			repoURL:       "github.com/owner/repo/",
			expectedPath:  "https://github.com/owner/repo",
			expectedError: false,
		},
		{
			repoURL:       "url",
			expectedPath:  "https://url/",
			expectedError: false,
		},
		{
			repoURL:       "github.com/owner/",
			expectedPath:  "",
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.repoURL, func(t *testing.T) {
			actualPath, err := ExtractRepoUrl(test.repoURL)
			if (err != nil) != test.expectedError {
				t.Errorf("expected error: %v, got: %v", test.expectedError, err)
			}
			if actualPath != test.expectedPath {
				t.Errorf("expected path: %v, got: %v", test.expectedPath, actualPath)
			}
		})
	}
}
