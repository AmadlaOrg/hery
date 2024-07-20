package validation

import (
	"github.com/AmadlaOrg/hery/entity/version"
)

// NewEntityVersionValidationService to set up the Entity version validation service
func NewEntityVersionValidationService() *VersionValidation {
	return &VersionValidation{
		Version: *version.NewEntityVersionService(),
	}
}
