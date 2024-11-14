package collection

import "github.com/AmadlaOrg/hery/storage"

// NewCollectionService to set up the collection service
func NewCollectionService() ICollection {
	return &SCollection{
		Storage: storage.NewStorageService(),
	}
}
