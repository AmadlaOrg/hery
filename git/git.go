package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

// Git interface to help with mocking
type Git interface {
	FetchRepo(url, dest string) error
	Tags(repoPath string) ([]string, error)
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
func Tags(repoPath string) ([]string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	refs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	var tags []string
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, ref.Name().Short())
		return nil
	})
	if err != nil {
		return nil, err
	}

	return tags, nil
}
