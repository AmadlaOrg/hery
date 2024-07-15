package validation

import (
	"github.com/AmadlaOrg/hery/entity/version"
	"github.com/AmadlaOrg/hery/entity/version/validation"
)

// NewEntityValidationService to set up the Entity Validation service
func NewEntityValidationService() *Validation {
	return &Validation{
		version.NewEntityVersionService(),
		validation.NewEntityVersionValidationService(),
	}
}
