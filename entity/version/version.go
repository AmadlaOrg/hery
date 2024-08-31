package version

import (
	"fmt"
	"github.com/AmadlaOrg/hery/util/git/remote"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Manager is an interface for managing versions.
type Manager interface {
	Extract(url string) (string, error)
	List(entityUrlPath string) ([]string, error)
	Latest(versions []string) (string, error)
	GeneratePseudo(entityFullRepoUrl string) (string, error)
}

// Service struct implements the Manager interface.
type Service struct {
	GitRemote remote.RepoRemoteManager
}

// Extract extracts the version from a go get URI.
func (v *Service) Extract(url string) (string, error) {
	versionAnnotationCount := strings.Count(url, "@")
	if versionAnnotationCount > 1 {
		return "", fmt.Errorf("invalid URI, contains more than one '@': %s", url)
	} else if versionAnnotationCount == 0 {
		return "", fmt.Errorf("no version found in URI: %s", url)
	}

	re := regexp.MustCompile(`@([^@]+)$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("no version found in URI: %s", url)
	}

	if !regexp.MustCompile(Format).MatchString(matches[1]) {
		return "", fmt.Errorf("invalid version format: %s", matches[1])
	}

	return matches[1], nil
}

// List returns a list of all the versions in tags with the format `v1.0.0`, `v1.0`, or `v1`.
func (v *Service) List(entityUrlPath string) ([]string, error) {
	tags, err := v.GitRemote.Tags(entityUrlPath)
	if err != nil {
		return nil, fmt.Errorf("error getting tags: %v", err)
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
func (v *Service) Latest(versions []string) (string, error) {
	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found")
	}

	sort.Slice(versions, func(i, j int) bool {
		return v.versionLess(versions[i], versions[j])
	})

	return versions[len(versions)-1], nil
}

// versionLess compares two version strings and returns true if v1 < v2.
func (v *Service) versionLess(v1, v2 string) bool {
	return v.compareVersions(v1, v2) < 0
}

// compareVersions compares two version strings and returns -1, 0, or 1 if v1 < v2, v1 == v2, or v1 > v2.
func (v *Service) compareVersions(v1, v2 string) int {
	parts1, pre1 := v.parseVersion(v1)
	parts2, pre2 := v.parseVersion(v2)

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
		return v.comparePreRelease(pre1, pre2)
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
func (v *Service) comparePreRelease(pre1, pre2 string) int {
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
func (v *Service) parseVersion(version string) ([]int, string) {
	re := regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)(?:-(alpha|beta|rc)\.(\d+))?$`)
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
		pre = matches[4] + "." + matches[5] // Ensures format like "beta.2"
	}

	return nums, pre
}

// GeneratePseudo generates a pseudo version to be used when there is no other source to identify the version of the entity.
func (v *Service) GeneratePseudo(entityFullRepoUrl string) (string, error) {
	commitHeadHash, err := v.GitRemote.CommitHeadHash(entityFullRepoUrl)
	if err != nil {
		return "", err
	}
	timestamp := time.Now().Format("20060102150405")
	pseudoVersion := fmt.Sprintf("v0.0.0-%s-%s", timestamp, commitHeadHash[:12])
	return pseudoVersion, nil
}
