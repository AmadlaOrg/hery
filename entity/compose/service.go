package compose

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/storage"
)

// New creates a new compose service.
func New(gitCfg *gitConfig.Config) Composer {
	return &composer{
		Storage: storage.New(),
		Entity:  entity.New(gitCfg),
	}
}
