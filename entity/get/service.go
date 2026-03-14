package get

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/build"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
)

// New to set up the Get service
func New(gitConfig *gitConfig.Config) Getter {
	return &getter{
		Entity:                  entity.New(gitConfig),
		EntityValidation:        validation.New(gitConfig),
		EntityVersion:           version.New(gitConfig),
		EntityVersionValidation: versionValidationPkg.New(gitConfig),
		Build:                   build.New(gitConfig),

		// Config
		GitConfig: gitConfig,
	}
}
