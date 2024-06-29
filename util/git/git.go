package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"os"
)

// Git interface to help with mocking
type Git interface {
	FetchRepo(url, dest string) error
	Tags(repoPath string) ([]string, error)
	CommitHeadHash(repoPath string) (string, error)
}

// FetchRepo clones the repository from the given URL to the specified destination.
func FetchRepo(url, dest string) error {
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

// Tags returns a list of tags for the repository at the specified path.
func Tags(repoPath string) (storer.ReferenceIter, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	tagRefs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	return tagRefs, nil
}

// CommitHeadHash retrieves the hash of the most recent commit
func CommitHeadHash(repoPath string) (string, error) {
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
