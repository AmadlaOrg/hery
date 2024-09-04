package remote

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

// IRemote to help with mocking
type IRemote interface {
	Tags(url string) ([]string, error)
	CommitHeadHash(url string) (string, error)
}

type SRemote struct{}

var (
	gitNewRemote     = git.NewRemote
	memoryNewStorage = memory.NewStorage
)

// Tags returns a list of tags for the repository at the specified URL.
func (s *SRemote) Tags(url string) ([]string, error) {
	rem := gitNewRemote(memoryNewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	refs, err := rem.List(&git.ListOptions{
		PeelingOption: git.IgnorePeeled,
	})
	if err != nil {
		return nil, err
	}

	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}

	return tags, nil
}

// CommitHeadHash retrieves the hash of the most recent commit
func (s *SRemote) CommitHeadHash(url string) (string, error) {
	rem := gitNewRemote(memoryNewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	// List all references from the remote repository
	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}

	var headRef *plumbing.Reference

	// Find the HEAD reference
	for _, ref := range refs {
		if ref.Name() == plumbing.HEAD {
			headRef = ref
			break
		}
	}

	if headRef == nil {
		return "", fmt.Errorf("HEAD reference not found")
	}

	var commitHash plumbing.Hash

	// Resolve symbolic reference if HEAD is symbolic
	if headRef.Type() == plumbing.SymbolicReference {
		// Find the reference that HEAD points to
		for _, ref := range refs {
			if ref.Name() == headRef.Target() {
				commitHash = ref.Hash()
				break
			}
		}
	} else {
		commitHash = headRef.Hash()
	}

	if commitHash.IsZero() {
		return "", fmt.Errorf("commit hash not found")
	}

	return commitHash.String(), nil
}
