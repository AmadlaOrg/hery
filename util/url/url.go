package url

import (
	"fmt"
	"net/url"
	"strings"
)

// EntityPathUrl removes the version in the Entity URI
func EntityPathUrl(entityUrl, version string) string {
	return strings.Replace(entityUrl, fmt.Sprintf("@%s", version), "", 1)
}

// EntityFullRepoUrl adds the protocol for a full URL to the repository in Git
func EntityFullRepoUrl(entityUrl string) string {
	return fmt.Sprintf("https://%s", entityUrl)
}

// ExtractRepoPath
func ExtractRepoPath(repoURL string) (string, error) {
	// Ensure the URL starts with a scheme for proper parsing
	if !strings.HasPrefix(repoURL, "http://") && !strings.HasPrefix(repoURL, "https://") {
		repoURL = "https://" + repoURL
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}

	// Trim the leading and trailing slashes from the path
	repoPath := strings.Trim(u.Path, "/")

	// Split the path into parts
	parts := strings.Split(repoPath, "/")

	if u.Host == "github.com" {
		// Check if the path is valid
		if len(parts) < 2 {
			return "", fmt.Errorf("invalid repository URL")
		}

		// Extract the repository path (owner/repo)
		repoPath = strings.Join(parts[:2], "/")
	}

	return "https://" + u.Host + "/" + repoPath, nil
}
