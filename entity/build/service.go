package build

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity"
	entityValidation "github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	entityVersionValidation "github.com/AmadlaOrg/hery/entity/version/validation"
)

// New to set up the entity Build service
func New(gitConfig *gitConfig.Config) Builder {
	return &builder{
		Entity:                  entity.New(gitConfig),
		EntityValidation:        entityValidation.New(gitConfig),
		EntityVersion:           version.New(gitConfig),
		EntityVersionValidation: entityVersionValidation.New(gitConfig),
	}
}
