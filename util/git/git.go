package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

// IGit to help with mocking
type IGit interface {
	FetchRepo(url, dest string) error
	CommitHeadHash(repoPath string) (string, error)
	CheckoutTag(repoPath, tagName string) error
}

type SGit struct{}

var (
	gitPlainOpen  = git.PlainOpen
	gitPlainClone = git.PlainClone
)

// FetchRepo clones the repository from the given URL to the specified destination.
func (s *SGit) FetchRepo(url, dest string) error {
	_, err := gitPlainClone(dest, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		// TODO: Add support for authentication because some people might require it
		// If you need authentication, add it here
		// Auth: &http.BasicAuth{
		//     Username: "your-username", // yes, this can be anything except an empty string
		//     Password: "your-token",
		// },
	})
	return err
}

// CommitHeadHash retrieves the hash of the most recent commit
func (s *SGit) CommitHeadHash(repoPath string) (string, error) {
	repo, err := gitPlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return "", err
	}

	// Get the commit object
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	return commit.Hash.String(), nil
}

// CheckoutTag checks out the specified branch or tag in the repository.
func (s *SGit) CheckoutTag(repoPath, tagName string) error {
	repo, err := gitPlainOpen(repoPath)
	if err != nil {
		return err
	}

	// Get the working tree
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Attempt to check out the reference as a branch
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", tagName)), //plumbing.NewBranchReferenceName(refName),
	})
	if err != nil {
		return err
	}

	return nil
}
