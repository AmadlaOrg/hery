package version

import (
	"fmt"
	"github.com/AmadlaOrg/hery/git"
	"os/exec"
	"regexp"
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

// List returns a list of all the versions in tags with the format `v1.0.0` or `v1.0` or `v1`
func List(dest string) ([]string, error) {
	tags, err := git.Tags(dest)
	if err != nil {
		return "", fmt.Errorf("error getting tags: %v\n", err)
	}

	for _, tag := range tags {

	}
}

// Latest
func Latest(versions []string) (string, error) {
	// TODO: Is it really getting the more recent?
	// TODO: Maybe we should filter so that it follows the format v1.0.0
	if len(versions) > 0 {
		return versions[len(versions)-1], nil // return the latest tag
	}

	return "", nil
}

// GeneratePseudo version to be used when there is no other source to identify the version of the entity
func GeneratePseudo(repoPath string) string {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	commitHash := strings.TrimSpace(string(output))
	timestamp := time.Now().Format("20060102150405")
	pseudoVersion := fmt.Sprintf("v0.0.0-%s-%s", timestamp, commitHash[:12])
	return pseudoVersion
}
