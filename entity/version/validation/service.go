package validation

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity/version"
)

// New to set up the Entity version validation service
func New(gitConfig *gitConfig.Config) Validator {
	return &validator{
		Version: version.New(gitConfig),
	}
}
