package version

import (
	"errors"
	"fmt"
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/LibraryUtils/git/remote"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// IVersion is an interface for managing versions.
type IVersion interface {
	Extract(url string) (string, error)
	List(entityUrlPath string) ([]string, error)
	Latest(versions []string) (string, error)
	GeneratePseudo(entityFullRepoUrl string) (string, error)
}

// SVersion struct implements the Manager interface.
type SVersion struct {
	GitRemoteConfig *gitConfig.Config
}

var (
	remoteNewGitRemoteService = remote.NewGitRemoteService
)

// Extract extracts the version from a go get URI.
//
// There are three types of versions that can be extracted:
// 1. Normal version (e.g.: v1.0.1, vx.x.x-beta.2)
// 2. Pseudo version (e.g.: v0.0.0-20240924093300-abcd1234efgh)
// 3. latest version (set or not set after the `@`)
//
// Normal version is detailed in the README at the root of the project and https://go.dev/doc/modules/version-numbers
// More information can also be found in: ./docs/version.md
//
// Pseudo version:
// - The first part will always contain: v0.0.0 // TODO: Golang seems the contain numbers at times... Maybe after a new commit following the most recent tagging of a version
// - The second part after the dash `-` is the date of when the entity was downloaded (for this reason it is "dynamic" and changing from one host to another) // TODO: Maybe this needs to be changed for the caching to be portable with pseudo versions. Second way: check the hash in the caching and then use the dating in the caching when it is re-downloaded (but this would make caching more of a higher authority)
// - The third part after the second dash `-` is a portion of the hash of the git commit that was downloaded (static and portable)
//
// The version type `latest` in two different ways:
// - If no version found the caller needs to check for the error type `ErrorExtractNoVersionFound` and determine if to set the entity version to `latest` or not
// - If it was set to then out put it as such
func (s *SVersion) Extract(url string) (string, error) {
	versionAnnotationCount := strings.Count(url, "@")
	if versionAnnotationCount > 1 {
		return "", errors.Join(ErrorExtractVersionAnnotationCountMoreThanOne, fmt.Errorf("url: %s", url))
	} else if versionAnnotationCount == 0 {
		return "", errors.Join(ErrorExtractNoVersionFound, fmt.Errorf("url: %s", url))
	}

	re := regexp.MustCompile(`@([^@]+)$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", errors.Join(ErrorExtractNoVersionFound, fmt.Errorf("url: %s", url))
	}

	if matches[1] == "latest" {
		return "latest", nil
	} else if !regexp.MustCompile(Format).MatchString(matches[1]) {
		return "", errors.Join(ErrorExtractInvalidVersion, fmt.Errorf("invalid version format: %s", matches[1]))
	}

	return matches[1], nil
}

// List returns a list of all the versions in tags with the format `v1.0.0`, `v1.0`, or `v1`.
func (s *SVersion) List(entityUrlPath string) ([]string, error) {
	gitRemote := remoteNewGitRemoteService(entityUrlPath, s.GitRemoteConfig)

	tags, err := gitRemote.Tags()
	if err != nil {
		return nil, errors.Join(ErrorListGitRemoteTags, err)
	}

	if len(tags) == 0 {
		return []string{}, nil
	}

	// Regular expression for matching version tags
	re := regexp.MustCompile(Format)

	var versions []string

	// Iterate over the tags and filter by the regex
	for _, tag := range tags {
		if re.MatchString(tag) {
			versions = append(versions, tag)
		}
	}

	return versions, nil
}

// Latest returns the most recent version from the list of versions.
func (s *SVersion) Latest(versions []string) (string, error) {
	if len(versions) == 0 {
		return "", ErrorLatestVersionsLenIsZero
	}

	sort.Slice(versions, func(i, j int) bool {
		return s.versionLess(versions[i], versions[j])
	})

	return versions[len(versions)-1], nil
}

// GeneratePseudo generates a pseudo version to be used when there is no other source to identify the version of the entity.
func (s *SVersion) GeneratePseudo(entityFullRepoUrl string) (string, error) {
	gitRemote := remoteNewGitRemoteService(entityFullRepoUrl, s.GitRemoteConfig)

	commitHeadHash, err := gitRemote.CommitHeadHash()
	if err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102150405")
	pseudoVersion := fmt.Sprintf("v0.0.0-%s-%s", timestamp, commitHeadHash[:12])

	return pseudoVersion, nil
}

//
// Private functions
//

// versionLess compares two version strings and returns true if v1 < v2.
func (s *SVersion) versionLess(v1, v2 string) bool {
	return s.compareVersions(v1, v2) < 0
}

// compareVersions compares two version strings and returns -1, 0, or 1 if v1 < v2, v1 == v2, or v1 > v2.
func (s *SVersion) compareVersions(v1, v2 string) int {
	parts1, pre1 := s.parseVersion(v1)
	parts2, pre2 := s.parseVersion(v2)

	// Compare major, minor, patch
	for i := 0; i < len(parts1); i++ {
		if parts1[i] < parts2[i] {
			return -1
		} else if parts1[i] > parts2[i] {
			return 1
		}
	}

	// Compare pre-release versions if both exist
	if pre1 != "" && pre2 != "" {
		return s.comparePreRelease(pre1, pre2)
	}

	// If only one has a pre-release version, it should be considered less
	if pre1 != "" {
		return -1
	}
	if pre2 != "" {
		return 1
	}

	// Versions are identical
	return 0
}

// comparePreRelease compares pre-release versions and returns -1, 0, or 1 if pre1 < pre2, pre1 == pre2, or pre1 > pre2.
func (s *SVersion) comparePreRelease(pre1, pre2 string) int {
	preOrder := map[string]int{"alpha": 0, "beta": 1, "rc": 2}

	parts1 := strings.Split(pre1, ".")
	parts2 := strings.Split(pre2, ".")

	// Compare pre-release type (alpha, beta, rc)
	order1 := preOrder[parts1[0]]
	order2 := preOrder[parts2[0]]

	if order1 < order2 {
		return -1
	} else if order1 > order2 {
		return 1
	}

	// Compare numeric suffix
	if len(parts1) > 1 && len(parts2) > 1 {
		num1, err1 := strconv.Atoi(parts1[1])
		num2, err2 := strconv.Atoi(parts2[1])

		if err1 == nil && err2 == nil {
			if num1 < num2 {
				return -1
			} else if num1 > num2 {
				return 1
			}
		}
	}

	// If no numeric suffix or they are the same
	return 0
}

// parseVersion parses a version string into its components and a pre-release identifier.
func (s *SVersion) parseVersion(version string) ([]int, string) {
	re := regexp.MustCompile(ParseVersionFormat)
	matches := re.FindStringSubmatch(version)

	if len(matches) == 0 {
		return nil, ""
	}

	nums := make([]int, 3)
	for i := 0; i < 3; i++ {
		var num int
		_, err := fmt.Sscanf(matches[i+1], "%d", &num)
		if err != nil {
			return nil, ""
		}
		nums[i] = num
	}

	pre := ""
	if len(matches[4]) > 0 && len(matches[5]) > 0 {
		pre = fmt.Sprintf("%s.%s", matches[4], matches[5]) // Ensures format like "beta.2"
	}

	return nums, pre
}
