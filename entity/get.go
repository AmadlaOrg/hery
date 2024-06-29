package entity

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	"github.com/AmadlaOrg/hery/util/git"
	"sync"
)

// Get
func Get(entityUrls []string, dest string) error {
	var wg sync.WaitGroup
	wg.Add(len(entityUrls))

	errCh := make(chan error, len(entityUrls)) // Channel to collect errors
	//TODO: entityPaths := make([]string, len(entityUrls)) Maybe check with directories

	for i, entityUrl := range entityUrls {
		go func(i int, entityUrl string) {
			defer wg.Done()

			// TODO: skip (continue) loop-iteration if entity with same version was already downloaded/installed maybe make a Map

			// Verify entity URL and build URL
			if !validation.EntityUrl(entityUrl) {
				errCh <- fmt.Errorf("invalid entity url: %s", entityUrl)
				return
			}
			url := fmt.Sprintf("https://%s", entityUrl)

			destination := fmt.Sprintf("/tmp/repo%d", i)

			// Download the Entity with `git clone`
			if err := git.FetchRepo(url, destination); err != nil {
				errCh <- fmt.Errorf("error fetching repo: %v\n", err)
			}

			// Extract or fetch the latest version
			entityVersion, err := version.Extract(entityUrl)
			// TODO: Make it global so that version.List(...) doesn't have to be called for the same Entity if look at again since you can have multiple versions of an entity
			entityVersionList, err := version.List(destination)
			if err != nil {
				entityVersion, err = version.Latest(entityVersionList)
				if err != nil || entityVersion == "" {
					entityVersion = version.GeneratePseudo(destination)
				} else {
					// TODO:
				}
			}

			//entityUrls[i].Version = entityVersion
		}(i, entityUrl)
	}

	wg.Wait()
	close(errCh) // Close the error channel

	var err error
	for e := range errCh {
		if err == nil {
			err = e
		} else {
			err = fmt.Errorf("%v; %v", err, e) // Combine multiple errors
		}
	}

	return err
}
