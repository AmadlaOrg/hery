package entity

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
)

// NewEntityService to set up the entity build service
func NewEntityService(gitConfig *gitConfig.Config) IEntity {
	return &SEntity{
		EntityVersion:           version.NewEntityVersionService(gitConfig),
		EntityVersionValidation: versionValidationPkg.NewEntityVersionValidationService(gitConfig),
		EntityValidation:        validation.NewEntityValidationService(gitConfig),
	}
}
