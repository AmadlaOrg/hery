package version

import "github.com/AmadlaOrg/hery/util/git/remote"

// NewEntityVersionService to set up the Entity Version Remote service
func NewEntityVersionService() *Version {
	return &Version{
		GitRemote: remote.NewGitRemoteService(),
	}
}
