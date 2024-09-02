package version

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/hery/util/git/remote"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExtract(t *testing.T) {
	// Initialize the service
	service := &Service{}

	// Scenario 1: Valid URL with a version
	t.Run("Valid URL with a version", func(t *testing.T) {
		url := "https://example.com/repo.git@v1.0.0"
		version, err := service.Extract(url)
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0", version)
	})

	// Scenario 2: URL with no version
	t.Run("URL with no version", func(t *testing.T) {
		url := "https://example.com/repo.git"
		version, err := service.Extract(url)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrorExtractNoVersionFound)
		assert.Equal(t, "", version)
	})

	// Scenario 3: URL with complex version string
	t.Run("URL with complex version string", func(t *testing.T) {
		url := "https://example.com/repo.git@v1.2.3-beta.1"
		version, err := service.Extract(url)
		assert.NoError(t, err)
		assert.Equal(t, "v1.2.3-beta.1", version)
	})

	// Scenario 4: URL with multiple '@' characters
	t.Run("URL with multiple '@' characters", func(t *testing.T) {
		url := "https://example.com/repo@branch@v1.0.0"
		version, err := service.Extract(url)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrorExtractVersionAnnotationCountMoreThanOne)
		assert.Empty(t, version)
	})

	// Scenario 5: URL with version not at the end
	t.Run("URL with version not at the end", func(t *testing.T) {
		url := "https://example.com/repo@v1.0.0/extra"
		version, err := service.Extract(url)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid version format:")
		assert.Empty(t, version)
	})
}

func TestList(t *testing.T) {
	// Scenario 1: Successful retrieval of tags
	t.Run("Successful retrieval of tags", func(t *testing.T) {
		// Create a new instance of the mocked GitRemote
		mockGitRemote := new(remote.MockGitRemote)

		// Create the service with the mocked GitRemote
		entityVersionService := &Service{
			GitRemote: mockGitRemote,
		}

		// Define the URL path to be used in the test
		entityUrlPath := "https://example.com/repo.git"

		// Define the expected tags and the result we want from the mock
		expectedTags := []string{"v1.0.0", "v2.0.0", "not-a-version", "v1"}
		mockGitRemote.On("Tags", entityUrlPath).Return(expectedTags, nil)

		// Call the List method
		result, err := entityVersionService.List(entityUrlPath)

		// Assert that there was no error
		assert.NoError(t, err)

		// Define the expected result after filtering
		expectedVersions := []string{"v1.0.0", "v2.0.0"}

		// Assert that the result matches the expected versions
		assert.Equal(t, expectedVersions, result)

		// Ensure that all expectations are met
		mockGitRemote.AssertExpectations(t)
	})

	// Scenario 2: Error scenario
	t.Run("Tags method returns an error", func(t *testing.T) {
		// Reinitialize the mock to clear the previous expectations
		mockGitRemote := new(remote.MockGitRemote)

		// Create the service with the mocked GitRemote
		entityVersionService := &Service{
			GitRemote: mockGitRemote,
		}

		// Define the URL path to be used in the test
		entityUrlPath := "https://example.com/repo.git"

		// Now let's test the case where the Tags method returns an error
		mockGitRemote.On("Tags", entityUrlPath).Return(nil, errors.New("git error"))

		// Call the List method again
		result, err := entityVersionService.List(entityUrlPath)

		// Assert that there was an error
		assert.Error(t, err)

		// Assert that the result is nil
		assert.Nil(t, result)

		// Ensure that all expectations are met
		mockGitRemote.AssertExpectations(t)
	})

	// Scenario 3: No tags found
	t.Run("No tags found", func(t *testing.T) {
		// Reinitialize the mock to clear the previous expectations
		mockGitRemote := new(remote.MockGitRemote)

		// Create the service with the mocked GitRemote
		entityVersionService := &Service{
			GitRemote: mockGitRemote,
		}

		// Define the URL path to be used in the test
		entityUrlPath := "https://example.com/repo.git"

		// Mock the scenario where the Tags method returns an empty list
		mockGitRemote.On("Tags", entityUrlPath).Return([]string{}, nil)

		// Call the List method again
		result, err := entityVersionService.List(entityUrlPath)

		// Assert that there was no error
		assert.NoError(t, err)

		// Assert that the result is an empty slice
		assert.Equal(t, []string{}, result)

		// Ensure that all expectations are met
		mockGitRemote.AssertExpectations(t)
	})
}

func TestLatest(t *testing.T) {
	// Initialize the service
	service := &Service{}

	// Scenario 1: No versions provided
	t.Run("No versions provided", func(t *testing.T) {
		var versions []string
		_, err := service.Latest(versions)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrorLatestVersionsLenIsZero)
	})

	// Scenario 2: Single version
	t.Run("Single version", func(t *testing.T) {
		versions := []string{"v1.0.0"}
		latest, err := service.Latest(versions)
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0", latest)
	})

	// Scenario 3: Multiple versions
	t.Run("Multiple versions", func(t *testing.T) {
		versions := []string{"v1.0.0", "v1.2.0", "v1.1.0"}
		latest, err := service.Latest(versions)
		assert.NoError(t, err)
		assert.Equal(t, "v1.2.0", latest)
	})

	// Scenario 4: Pre-release versions
	t.Run("Pre-release versions", func(t *testing.T) {
		versions := []string{"v1.0.0-alpha.1", "v1.0.0-beta.1", "v1.0.0"}
		latest, err := service.Latest(versions)
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0", latest)
	})

	t.Run("Pre-release versions", func(t *testing.T) {
		versions := []string{"v1.0.0-alpha.1", "v1.0.0-beta.1"}
		latest, err := service.Latest(versions)
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0-beta.1", latest)
	})

	// Scenario 5: Complex versions with pre-releases
	t.Run("Complex versions with pre-releases", func(t *testing.T) {
		versions := []string{"v1.0.0-beta.2", "v1.0.0-alpha.1", "v1.0.0-beta.1", "v1.0.0"}
		latest, err := service.Latest(versions)
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0", latest)
	})

	t.Run("Complex versions with pre-releases", func(t *testing.T) {
		versions := []string{"v1.0.0-beta.2", "v1.0.0-alpha.1", "v1.0.0-beta.1"}
		latest, err := service.Latest(versions)
		assert.NoError(t, err)
		assert.Equal(t, "v1.0.0-beta.2", latest)
	})
}

func TestVersionLess(t *testing.T) {
	service := &Service{}

	// Scenario 1: Basic comparison
	t.Run("Basic comparison", func(t *testing.T) {
		assert.True(t, service.versionLess("v1.0.0", "v1.0.1"))
		assert.False(t, service.versionLess("v1.0.1", "v1.0.0"))
		assert.False(t, service.versionLess("v1.0.0", "v1.0.0"))
	})

	// Scenario 2: Major version difference
	t.Run("Major version difference", func(t *testing.T) {
		assert.True(t, service.versionLess("v1.0.0", "v2.0.0"))
		assert.False(t, service.versionLess("v2.0.0", "v1.0.0"))
	})

	// Scenario 3: Pre-release versions
	t.Run("Pre-release versions", func(t *testing.T) {
		assert.True(t, service.versionLess("v1.0.0-alpha.1", "v1.0.0"))
		assert.False(t, service.versionLess("v1.0.0", "v1.0.0-alpha.1"))
		assert.True(t, service.versionLess("v1.0.0-alpha.1", "v1.0.0-beta.1"))
	})
}

func TestCompareVersions(t *testing.T) {
	service := &Service{}

	// Scenario 1: Equal versions
	t.Run("Equal versions", func(t *testing.T) {
		assert.Equal(t, 0, service.compareVersions("v1.0.0", "v1.0.0"))
	})

	// Scenario 2: Different versions
	t.Run("Different versions", func(t *testing.T) {
		assert.Equal(t, -1, service.compareVersions("v1.0.0", "v1.1.0"))
		assert.Equal(t, 1, service.compareVersions("v1.2.0", "v1.1.0"))
	})

	// Scenario 3: Pre-release comparison
	t.Run("Pre-release comparison", func(t *testing.T) {
		assert.Equal(t, -1, service.compareVersions("v1.0.0-alpha.1", "v1.0.0"))
		assert.Equal(t, 1, service.compareVersions("v1.0.0", "v1.0.0-alpha.1"))
		assert.Equal(t, -1, service.compareVersions("v1.0.0-alpha.1", "v1.0.0-beta.1"))
	})
}

func TestComparePreRelease(t *testing.T) {
	service := &Service{}

	// Scenario 1: Equal pre-release versions
	t.Run("Equal pre-release versions", func(t *testing.T) {
		assert.Equal(t, 0, service.comparePreRelease("alpha.1", "alpha.1"))
	})

	// Scenario 2: Different pre-release versions
	t.Run("Different pre-release versions", func(t *testing.T) {
		assert.Equal(t, -1, service.comparePreRelease("alpha.1", "beta.1"))
		assert.Equal(t, 1, service.comparePreRelease("beta.2", "alpha.1"))
	})

	// Scenario 3: Numeric suffix comparison
	t.Run("Numeric suffix comparison", func(t *testing.T) {
		assert.Equal(t, -1, service.comparePreRelease("alpha.1", "alpha.2"))
		assert.Equal(t, 1, service.comparePreRelease("alpha.3", "alpha.2"))
	})
}

func TestParseVersion(t *testing.T) {
	service := &Service{}

	// Scenario 1: Parse basic version
	t.Run("Parse basic version", func(t *testing.T) {
		parts, pre := service.parseVersion("v1.2.3")
		assert.Equal(t, []int{1, 2, 3}, parts)
		assert.Equal(t, "", pre)
	})

	// Scenario 2: Parse version with pre-release
	t.Run("Parse version with pre-release", func(t *testing.T) {
		parts, pre := service.parseVersion("v1.2.3-alpha.1")
		assert.Equal(t, []int{1, 2, 3}, parts)
		assert.Equal(t, "alpha.1", pre)
	})

	// Scenario 3: Parse incomplete version
	t.Run("Parse incomplete version", func(t *testing.T) {
		parts, pre := service.parseVersion("v1.0")
		assert.Nil(t, parts)
		assert.Empty(t, pre)
	})
}

func TestGeneratePseudo(t *testing.T) {
	// Scenario 1: Successful generation of pseudo-version
	t.Run("Successful generation of pseudo-version", func(t *testing.T) {
		// Create a new instance of the mocked GitRemote
		mockGitRemote := new(remote.MockGitRemote)

		// Create the service with the mocked GitRemote
		entityVersionService := &Service{
			GitRemote: mockGitRemote,
		}

		// Define the URL path to be used in the test
		entityFullRepoUrl := "https://example.com/repo.git"

		// Mock the CommitHeadHash method to return a specific hash
		expectedCommitHash := "abcd1234efgh5678ijkl"
		mockGitRemote.On("CommitHeadHash", entityFullRepoUrl).Return(expectedCommitHash, nil)

		// Call the GeneratePseudo method
		pseudoVersion, err := entityVersionService.GeneratePseudo(entityFullRepoUrl)

		// Assert that there was no error
		assert.NoError(t, err)

		// Build the expected pseudo version format
		expectedTimestamp := time.Now().Format("20060102150405")
		expectedPseudoVersion := fmt.Sprintf("v0.0.0-%s-%s", expectedTimestamp, expectedCommitHash[:12])

		// Assert that the pseudo version matches the expected format
		assert.Equal(t, expectedPseudoVersion, pseudoVersion)

		// Ensure that all expectations are met
		mockGitRemote.AssertExpectations(t)
	})

	// Scenario 2: Error case when retrieving commit hash
	t.Run("Error case when retrieving commit hash", func(t *testing.T) {
		// Reinitialize the mock to clear the previous expectations
		mockGitRemote := new(remote.MockGitRemote)

		// Create the service with the mocked GitRemote
		entityVersionService := &Service{
			GitRemote: mockGitRemote,
		}

		// Define the URL path to be used in the test
		entityFullRepoUrl := "https://example.com/repo.git"

		// Mock the CommitHeadHash method to return an error
		mockGitRemote.On("CommitHeadHash", entityFullRepoUrl).Return("", errors.New("commit hash error"))

		// Call the GeneratePseudo method
		pseudoVersion, err := entityVersionService.GeneratePseudo(entityFullRepoUrl)

		// Assert that there was an error
		assert.Error(t, err)

		// Assert that the pseudo version is empty due to the error
		assert.Equal(t, "", pseudoVersion)

		// Ensure that all expectations are met
		mockGitRemote.AssertExpectations(t)
	})
}
