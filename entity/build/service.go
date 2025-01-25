package build

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity"
	entityValidation "github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	entityVersionValidation "github.com/AmadlaOrg/hery/entity/version/validation"
)

// NewEntityBuildService to set up the entity Build service
func NewEntityBuildService(gitConfig *gitConfig.Config) IBuild {
	return &SBuild{
		Entity:                  entity.NewEntityService(gitConfig),
		EntityValidation:        entityValidation.NewEntityValidationService(gitConfig),
		EntityVersion:           version.NewEntityVersionService(gitConfig),
		EntityVersionValidation: entityVersionValidation.NewEntityVersionValidationService(gitConfig),
	}
}
