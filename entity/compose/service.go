package compose

import "github.com/AmadlaOrg/hery/storage"

// NewComposeService to set up the compose service
func NewComposeService() *Composer {
	return &Composer{
		Storage: storage.NewStorageService(),
	}
}
