package validation

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity/schema"
	schemaValidationPkg "github.com/AmadlaOrg/hery/entity/schema/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	"github.com/AmadlaOrg/hery/entity/version/validation"
)

// New to set up the Entity Validation service
func New(gitConfig *gitConfig.Config) Validator {
	return &validator{
		Version:           version.New(gitConfig),
		VersionValidation: validation.New(gitConfig),
		Schema:            schema.New(),
		SchemaValidation:  schemaValidationPkg.New(),
	}
}
