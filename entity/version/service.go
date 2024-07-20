package version

import "github.com/AmadlaOrg/hery/util/git/remote"

// NewEntityVersionService to set up the Entity Version Remote service
func NewEntityVersionService() *Service {
	return &Service{
		GitRemote: remote.NewGitRemoteService(),
	}
}
