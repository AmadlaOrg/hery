package url

import (
	"fmt"
	"strings"
)

func EntityPathUrl(entityUrl, version string) string {
	return strings.Replace(entityUrl, fmt.Sprintf("@%s", version), "", 1)
}

func EntityFullRepoUrl(entityUrl string) string {
	return fmt.Sprintf("https://%s", entityUrl)
}
