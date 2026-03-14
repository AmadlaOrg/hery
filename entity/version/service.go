package version

import gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"

// New to set up the Entity Version Remote service
func New(gitConfig *gitConfig.Config) Version {
	return &versionImpl{
		GitRemoteConfig: gitConfig,
	}
}
