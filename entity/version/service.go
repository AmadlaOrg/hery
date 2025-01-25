package version

import gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"

// NewEntityVersionService to set up the Entity Version Remote service
func NewEntityVersionService(gitConfig *gitConfig.Config) IVersion {
	return &SVersion{
		GitRemoteConfig: gitConfig,
	}
}
