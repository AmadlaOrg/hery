package validation

import (
	"github.com/AmadlaOrg/hery/entity/schema"
	schemaValidationPkg "github.com/AmadlaOrg/hery/entity/schema/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	"github.com/AmadlaOrg/hery/entity/version/validation"
)

// NewEntityValidationService to set up the Entity Validation service
func NewEntityValidationService() *SValidation {
	return &SValidation{
		Version:           version.NewEntityVersionService(),
		VersionValidation: validation.NewEntityVersionValidationService(),
		Schema:            schema.NewEntitySchemaService(),
		SchemaValidation:  schemaValidationPkg.NewEntitySchemaValidationService(),
	}
}
