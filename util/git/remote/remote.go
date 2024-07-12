package remote

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"log"
)

// Interface to help with mocking
type Interface interface {
	Tags(repoPath string) ([]string, error)
	CommitHeadHash(repoPath string) (string, error)
}

type GitRemote struct{}

// Tags returns a list of tags for the repository at the specified path.
func (gt *GitRemote) Tags(url string) ([]string, error) {
	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	// We can then use every Remote functions to retrieve wanted information
	refs, err := rem.List(&git.ListOptions{
		// Returns all references, including peeled references.
		PeelingOption: git.AppendPeeled,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Filters the references list and only keeps tags
	var tags []string
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}

	return tags, nil
}

// CommitHeadHash retrieves the hash of the most recent commit
func (gt *GitRemote) CommitHeadHash(url string) (string, error) {
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
