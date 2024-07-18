package compose

import "github.com/AmadlaOrg/hery/storage"

// NewComposeService to set up the compose service
func NewComposeService() *Compose {
	return &Compose{
		Storage: storage.NewStorageService(),
	}
}
