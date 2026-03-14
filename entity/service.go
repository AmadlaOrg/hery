package entity

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
)

// New to set up the entity build service
func New(gitConfig *gitConfig.Config) Service {
	return &service{
		EntityVersion:           version.New(gitConfig),
		EntityVersionValidation: versionValidationPkg.New(gitConfig),
		EntityValidation:        validation.New(gitConfig),
	}
}
