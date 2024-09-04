package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func Test_Integration_FetchRepo(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "hery_test_*")
	if err != nil {
		t.Fatal(err)
	}

	// Clean up after test
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	tests := []struct {
		name        string
		url         string
		dest        string
		expectedErr bool
	}{
		{
			name:        "Successful clone",
			url:         "https://github.com/AmadlaOrg/QAFixturesEntityPseudoVersion",
			dest:        filepath.Join(tempDir, "qa"),
			expectedErr: false,
		},
		{
			name:        "Failed clone",
			url:         "github.com/AmadlaOrg/QANoneExistingRepo",
			dest:        filepath.Join(tempDir, "qa"),
			expectedErr: true,
		},
	}

	gitService := NewGitService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gitService.FetchRepo(tt.url, tt.dest)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Check if the .git directory exists in the destination
				gitDir := filepath.Join(tt.dest, ".git")
				_, err := os.Stat(gitDir)
				assert.NoError(t, err, "Expected .git directory to exist")
				assert.True(t, !os.IsNotExist(err), ".git directory does not exist")
			}
		})
	}
}

func Test_Integration_CommitHeadHash(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "hery_test_*")
	if err != nil {
		t.Fatal(err)
	}

	// Clean up after test
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	gitService := NewGitService()

	err = gitService.FetchRepo("https://github.com/AmadlaOrg/QAFixturesEntityPseudoVersion", tempDir)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := gitService.CommitHeadHash(tempDir)

	assert.NoError(t, err)
	assert.Equal(t, "a33efb99e6c7d182034a5c5c2cb7a165026bff84", hash)
}

// Mock the PlainOpen function to return an error
func Test_Integration_CommitHeadHash_RepoOpenError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "hery_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	// Ensure that the directory is removed to simulate a repo not existing
	err = os.RemoveAll(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	gitService := NewGitService()

	_, err = gitService.CommitHeadHash(tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository does not exist") // Customize this based on actual error message
}

func Test_Integration_CommitHeadHash_RepoHeadError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "hery_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	// Setup a minimal git repo
	_, err = git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatal(err)
	}

	gitService := NewGitService()

	// Write an invalid reference to the HEAD file
	headFilePath := filepath.Join(tempDir, ".git", "HEAD")
	err = os.WriteFile(headFilePath, []byte("ref: refs/heads/invalid-branch\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = gitService.CommitHeadHash(tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reference not found") // Customize this based on actual error message
}

// Mock the CommitObject function to return an error
func Test_Integration_CommitHeadHash_CommitObjectError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "hery_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	// Initialize the repository
	repo, _ := git.PlainInit(tempDir, false)

	// Manually create an invalid commit object
	ref := plumbing.NewHashReference(plumbing.HEAD, plumbing.NewHash("0000000000000000000000000000000000000000"))
	err = repo.Storer.SetReference(ref)
	if err != nil {
		t.Fatal(err)
	}

	gitService := NewGitService()

	_, err = gitService.CommitHeadHash(tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "object not found") // Customize this based on actual error message
}

// Test repository open error
func Test_Integration_CheckoutTag_RepoOpenError(t *testing.T) {
	// Use a non-existing directory to simulate repository open error
	tempDir := "/non/existing/path"

	gitService := NewGitService()
	err := gitService.CheckoutTag(tempDir, "v1.0.0")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repository does not exist") // Customize based on actual error message
}

// Test checkout error (e.g., tag does not exist)
func Test_Integration_CheckoutTag_CheckoutError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "hery_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}(tempDir)

	// Initialize a new git repository
	_, err = git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatal(err)
	}

	gitService := NewGitService()
	err = gitService.CheckoutTag(tempDir, "non-existing-tag")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reference not found") // Customize based on actual error message
}
