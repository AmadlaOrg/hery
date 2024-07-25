package url

import (
	"fmt"
	"net/url"
	"strings"
)

// TrimVersion removes the version in the Entity URI
func TrimVersion(entityUrl, version string) string {
	return strings.Replace(entityUrl, fmt.Sprintf("@%s", version), "", 1)
}

// ExtractRepoUrl adds https:// to the repository URL pass by param
//
// It also checks if the repoURL is of `github.com` and if it is it then does a specific validation on top of an extractions
// meaning if there are more path parts then it only returns the full URL for the repository itself.
func ExtractRepoUrl(repoURL string) (string, error) {
	if !strings.HasPrefix(repoURL, "http://") && !strings.HasPrefix(repoURL, "https://") {
		repoURL = fmt.Sprintf("https://%s", repoURL)
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

	return fmt.Sprintf("https://%s/%s", u.Host, repoPath), nil
}
