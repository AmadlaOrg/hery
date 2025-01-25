package validation

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity/schema"
	schemaValidationPkg "github.com/AmadlaOrg/hery/entity/schema/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	"github.com/AmadlaOrg/hery/entity/version/validation"
)

// NewEntityValidationService to set up the Entity Validation service
func NewEntityValidationService(gitConfig *gitConfig.Config) IValidation {
	return &SValidation{
		Version:           version.NewEntityVersionService(gitConfig),
		VersionValidation: validation.NewEntityVersionValidationService(gitConfig),
		Schema:            schema.NewEntitySchemaService(),
		SchemaValidation:  schemaValidationPkg.NewEntitySchemaValidationService(),
	}
}
