package build

import (
	"github.com/AmadlaOrg/hery/entity"
	entityValidation "github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	entityVersionValidation "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/util/git"
)

// NewEntityBuildService to set up the entity Build service
func NewEntityBuildService() *SBuild {
	return &SBuild{
		Git:                     git.NewGitService(),
		Entity:                  entity.NewEntityService(),
		EntityValidation:        entityValidation.NewEntityValidationService(),
		EntityVersion:           version.NewEntityVersionService(),
		EntityVersionValidation: entityVersionValidation.NewEntityVersionValidationService(),
	}
}
