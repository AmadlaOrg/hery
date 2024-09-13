package validation

import schemaPkg "github.com/AmadlaOrg/hery/entity/schema"

// NewEntitySchemaValidationService to set up the entity Validation service
func NewEntitySchemaValidationService() *SValidation {
	return &SValidation{
		Schema: schemaPkg.NewEntitySchemaService(),
	}
}
