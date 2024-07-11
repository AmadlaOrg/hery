package version

import (
	"fmt"
	"github.com/AmadlaOrg/hery/util/git/remote"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Extract extracts the version from a go get URI.
func Extract(url string) (string, error) {
	re := regexp.MustCompile(`@(.+)$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("no version found in URI: %s", url)
	}
	return matches[1], nil
}

// Exists checks if a specific version exists in the list of versions
func Exists(entityUrlPath, version string) (bool, error) {
	versions, err := List(entityUrlPath)
	if err != nil {
		return false, err
	}

	for _, v := range versions {
		if v == version {
			return true, nil
		}
	}

	return false, nil
}

// List returns a list of all the versions in tags with the format `v1.0.0` or `v1.0` or `v1`
func List(entityUrlPath string) ([]string, error) {
	tags, err := remote.Tags(entityUrlPath)
	if err != nil {
		return nil, fmt.Errorf("error getting tags: %v\n", err)
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

// Latest returns the most recent version from the list of versions
func Latest(versions []string) (string, error) {
	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found")
	}

	sort.Slice(versions, func(i, j int) bool {
		return versionLess(versions[i], versions[j])
	})

	return versions[len(versions)-1], nil
}

// versionLess compares two version strings and returns true if v1 < v2
func versionLess(v1, v2 string) bool {
	return compareVersions(v1, v2) < 0
}

// compareVersions compares two version strings and returns -1, 0, or 1 if v1 < v2, v1 == v2, or v1 > v2
func compareVersions(v1, v2 string) int {
	parts1, pre1 := parseVersion(v1)
	parts2, pre2 := parseVersion(v2)

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		if parts1[i] < parts2[i] {
			return -1
		} else if parts1[i] > parts2[i] {
			return 1
		}
	}

	// Compare lengths if one version has more parts
	if len(parts1) < len(parts2) {
		return -1
	} else if len(parts1) > len(parts2) {
		return 1
	}

	// Compare pre-release versions if they exist
	if pre1 == "" && pre2 != "" {
		return 1
	} else if pre1 != "" && pre2 == "" {
		return -1
	} else if pre1 != "" && pre2 != "" {
		return comparePreRelease(pre1, pre2)
	}

	return 0
}

// comparePreRelease compares pre-release versions and returns -1, 0, or 1 if pre1 < pre2, pre1 == pre2, or pre1 > pre2
func comparePreRelease(pre1, pre2 string) int {
	preOrder := map[string]int{"alpha": 0, "beta": 1, "rc": 2}
	parts1 := strings.Split(pre1, ".")
	parts2 := strings.Split(pre2, ".")

	order1 := preOrder[parts1[0]]
	order2 := preOrder[parts2[0]]

	if order1 < order2 {
		return -1
	} else if order1 > order2 {
		return 1
	}

	num1, num2 := 0, 0
	_, err := fmt.Sscanf(parts1[1], "%d", &num1)
	if err != nil {
		return 0
	}
	_, err = fmt.Sscanf(parts2[1], "%d", &num2)
	if err != nil {
		return 0
	}

	if num1 < num2 {
		return -1
	} else if num1 > num2 {
		return 1
	}

	return 0
}

// parseVersion parses a version string into its components and a pre-release identifier
func parseVersion(version string) ([]int, string) {
	re := regexp.MustCompile(Format)
	matches := re.FindStringSubmatch(version)

	var nums []int
	for i := 1; i <= 3; i++ {
		if matches[i] != "" {
			var num int
			_, err := fmt.Sscanf(strings.TrimPrefix(matches[i], "."), "%d", &num)
			if err != nil {
				return nil, ""
			}
			nums = append(nums, num)
		}
	}

	pre := ""
	if matches[4] != "" {
		pre = matches[4][1:] // Remove leading '-'
	}

	return nums, pre
}

// GeneratePseudo version to be used when there is no other source to identify the version of the entity
func GeneratePseudo(entityUrlPath string) (string, error) {
	commitHeadHash, err := remote.CommitHeadHash(entityUrlPath)
	if err != nil {
		return "", err
	}
	timestamp := time.Now().Format("20060102150405")
	pseudoVersion := fmt.Sprintf("v0.0.0-%s-%s", timestamp, commitHeadHash[:12])
	return pseudoVersion, nil
}
