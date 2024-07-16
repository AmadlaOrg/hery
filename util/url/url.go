package url

import (
	"fmt"
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
