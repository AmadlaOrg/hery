package remote

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

// RepoRemoteManager to help with mocking
type RepoRemoteManager interface {
	Tags(url string) ([]string, error)
	CommitHeadHash(url string) (string, error)
}

type GitRemote struct{}

// Tags returns a list of tags for the repository at the specified URL.
func (gr *GitRemote) Tags(url string) ([]string, error) {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	// Retrieve all references from the remote repository
	refs, err := rem.List(&git.ListOptions{
		// Returns all references, including peeled references.
		PeelingOption: git.AppendPeeled,
	})
	if err != nil {
		return nil, err
	}

	// Filter the references list and only keep tags
	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}

	return tags, nil
}

// CommitHeadHash retrieves the hash of the most recent commit
func (gr *GitRemote) CommitHeadHash(url string) (string, error) {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
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

	// Get the commit object for the HEAD reference
	commitHash := headRef.Hash()

	return commitHash.String(), nil
}
