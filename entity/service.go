package entity

import (
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
)

// NewEntityService to set up the entity build service
func NewEntityService() *SEntity {
	return &SEntity{
		EntityVersion:           version.NewEntityVersionService(),
		EntityVersionValidation: versionValidationPkg.NewEntityVersionValidationService(),

		// Data
		Entities: []Entity{},
	}
}
