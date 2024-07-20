package build

import (
	entityValidation "github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	entityVersionValidation "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/util/git"
)

// NewEntityBuildService to set up the entity build service
func NewEntityBuildService() *Builder {
	return &Builder{
		Git:                     git.NewGitService(),
		EntityValidation:        entityValidation.NewEntityValidationService(),
		EntityVersion:           version.NewEntityVersionService(),
		EntityVersionValidation: *entityVersionValidation.NewEntityVersionValidationService(),
	}
}
