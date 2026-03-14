package validation

import schemaPkg "github.com/AmadlaOrg/hery/entity/schema"

// New to set up the entity Validation service
func New() Validator {
	return &validator{
		Schema: schemaPkg.New(),
	}
}
