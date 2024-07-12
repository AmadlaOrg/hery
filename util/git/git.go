package git

import (
	"github.com/go-git/go-git/v5"
	"os"
)

// Interface to help with mocking
type Interface interface {
	FetchRepo(url, dest string) error
	CommitHeadHash(repoPath string) (string, error)
}

type Git struct{}

// FetchRepo clones the repository from the given URL to the specified destination.
func (g *Git) FetchRepo(url, dest string) error {
	_, err := git.PlainClone(dest, false, &git.CloneOptions{
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
func (g *Git) CommitHeadHash(repoPath string) (string, error) {
	// Open the repository
	repo, err := git.PlainOpen(repoPath)
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
