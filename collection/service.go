package collection

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
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
func NewEntityCollectionService(gitConfig *gitConfig.Config) IEntityCollection {
	return &SEntityCollection{
		EntityVersion:           version.NewEntityVersionService(gitConfig),
		EntityVersionValidation: versionValidationPkg.NewEntityVersionValidationService(gitConfig),
		EntityValidation:        validation.NewEntityValidationService(gitConfig),

		// Data
		//Entities: &[]*entity.Entity{},
		Collection: &Collection{},
	}
}
