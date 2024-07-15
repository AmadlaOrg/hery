package entity

import (
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	utilGit "github.com/AmadlaOrg/hery/util/git"
)

// NewGetService to set up the Get service
func NewGetService() *GetService {
	return &GetService{
		Git:                     utilGit.NewGitService(),
		EntityValidation:        validation.NewEntityValidationService(),
		EntityVersion:           version.NewEntityVersionService(),
		EntityVersionValidation: versionValidationPkg.NewEntityVersionValidationService(),
	}
}
