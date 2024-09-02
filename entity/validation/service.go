package validation

import (
	"github.com/AmadlaOrg/hery/entity/version"
	"github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/heryext/schema"
)

// NewEntityValidationService to set up the Entity Validation service
func NewEntityValidationService() *SValidation {
	return &SValidation{
		Version:           version.NewEntityVersionService(),
		VersionValidation: validation.NewEntityVersionValidationService(),
		Schema:            schema.NewSchemaService(),
	}
}
