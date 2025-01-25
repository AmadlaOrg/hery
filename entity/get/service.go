package get

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/build"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
)

// NewGetService to set up the Get service
func NewGetService(gitConfig *gitConfig.Config) IGet {
	return &SGet{
		Entity:                  entity.NewEntityService(gitConfig),
		EntityValidation:        validation.NewEntityValidationService(gitConfig),
		EntityVersion:           version.NewEntityVersionService(gitConfig),
		EntityVersionValidation: versionValidationPkg.NewEntityVersionValidationService(gitConfig),
		Build:                   build.NewEntityBuildService(gitConfig),

		// Config
		GitConfig: gitConfig,
	}
}
