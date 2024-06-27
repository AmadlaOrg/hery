package version

import (
	"fmt"
	"github.com/AmadlaOrg/hery/git"
	"os/exec"
	"strings"
	"time"
)

func Latest(url, dest string) string {
	if err := git.FetchRepo(url, dest); err != nil {
		fmt.Printf("Error fetching repo: %v\n", err)
		return ""
	}

	tags, err := git.Tags(dest)
	if err != nil {
		fmt.Printf("Error getting tags: %v\n", err)
		return ""
	}

	if len(tags) > 0 {
		return tags[len(tags)-1] // return the latest tag
	}

	pseudoVersion, err := generatePseudo(dest)
	if err != nil {
		fmt.Printf("Error generating pseudo-version: %v\n", err)
		return ""
	}

	return pseudoVersion
}

// generatePseudo version to be used when there is no other source to identify the version of the entity
func generatePseudo(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	commitHash := strings.TrimSpace(string(output))
	timestamp := time.Now().Format("20060102150405")
	pseudoVersion := fmt.Sprintf("v0.0.0-%s-%s", timestamp, commitHash[:12])
	return pseudoVersion, nil
}
