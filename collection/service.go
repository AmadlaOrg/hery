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
