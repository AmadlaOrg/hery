package collection

import (
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/storage"
)

// NewCollectionService to set up the collection service
func NewCollectionService() ICollection {
	return &SCollection{
		Storage: storage.NewStorageService(),
	}
}

// NewEntityCollectionService to set up the entity collection service
func NewEntityCollectionService() IEntityCollection {
	return &SEntityCollection{
		EntityVersion:           version.NewEntityVersionService(),
		EntityVersionValidation: versionValidationPkg.NewEntityVersionValidationService(),
		EntityValidation:        validation.NewEntityValidationService(),

		// Data
		//Entities: &[]*entity.Entity{},
		Collection: &Collection{},
	}
}
