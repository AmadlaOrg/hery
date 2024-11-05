package version

import "github.com/AmadlaOrg/hery/util/git/remote"

// NewEntityVersionService to set up the Entity Version Remote service
func NewEntityVersionService() IVersion {
	return &SVersion{
		GitRemote: remote.NewGitRemoteService(),
	}
}
