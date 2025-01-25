package validation

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity/version"
)

// NewEntityVersionValidationService to set up the Entity version validation service
func NewEntityVersionValidationService(gitConfig *gitConfig.Config) IValidation {
	return &SValidation{
		Version: version.NewEntityVersionService(gitConfig),
	}
}
