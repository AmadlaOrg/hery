package compose

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/storage"
)

// NewComposeService creates a new compose service.
func NewComposeService(gitCfg *gitConfig.Config) IComposer {
	return &SComposer{
		Storage: storage.NewStorageService(),
		Entity:  entity.NewEntityService(gitCfg),
	}
}
