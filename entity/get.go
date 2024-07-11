package entity

import (
	"fmt"
	"github.com/AmadlaOrg/hery/collection"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	utilFilePkg "github.com/AmadlaOrg/hery/util/file"
	"github.com/AmadlaOrg/hery/util/git"
	"github.com/AmadlaOrg/hery/util/url"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

// Get with collection name and the args that are the entities urls, calls on download to get the entities
func Get(collectionName, storagePath string, args []string) {
	// Validate that all the URLs pass in the arguments are valid
	if len(args) == 0 {
		log.Fatal("No entity URL(s) specified")
	}

	for _, arg := range args {
		if !validation.EntityUrl(arg) {
			log.Fatalf("Invalid entity URL: %s", arg)
		}
	}

	collectionStoragePath := collection.Path(collectionName, storagePath)
	if !utilFilePkg.Exists(collectionStoragePath) {
		log.Fatalf("The collection storage directory does not exist: %s", collectionStoragePath)
	}

	println(collectionStoragePath)
	err := download(args, collectionStoragePath)
	if err != nil {
		log.Fatal(err)
	}
}

// download in parallel all the entities
func download(entityUrls []string, collectionStoragePath string) error {
	var wg sync.WaitGroup
	wg.Add(len(entityUrls))

	errCh := make(chan error, len(entityUrls)) // Channel to collect errors
	//TODO: entityPaths := make([]string, len(entityUrls)) Maybe check with directories

	for i, entityUrl := range entityUrls {
		go func(i int, entityUrl string) {
			defer wg.Done()

			// TODO: skip (continue) loop-iteration if entity with same version was already downloaded/installed maybe make a Map

			var entityUrlPath string
			var entityFullRepoUrl string
			var entityVersion string
			if strings.Contains(entityUrl, "@") {
				entityVersion, versionExtractErr := version.Extract(entityUrl)
				entityUrlPath = url.EntityPathUrl(entityUrl, entityVersion)
				entityFullRepoUrl = url.EntityFullRepoUrl(entityUrlPath)
				entityVersionList, err := version.List(entityFullRepoUrl)

				// If no tags (versions) are found there will be no errors but there might be some error caused by git
				if err != nil {
					errCh <- err
				} else {

					// True if no version was extracted after @
					if versionExtractErr != nil {
						errCh <- versionExtractErr
					} else {
						// From the git repo there are no tags then it will be `0`
						if len(entityVersionList) == 0 {
							entityVersion, err = version.GeneratePseudo(entityFullRepoUrl)
						} else { // Tries to find the newest version from git repo tags
							entityVersion, err = version.Latest(entityVersionList)
						}
					}
				}
			} else {
				entityUrlPath = entityUrl
				entityFullRepoUrl = url.EntityFullRepoUrl(entityUrlPath)
			}

			destination := filepath.Join(collectionStoragePath, entityUrl)

			// Download the Entity with `git clone`
			if err := git.FetchRepo(entityFullRepoUrl, destination); err != nil {
				errCh <- fmt.Errorf("error fetching repo: %v\n", err)
			}

			// Extract or fetch the latest version
			/*entityVersion, err := version.Extract(entityUrl)
			if err != nil {
				//errCh <- err

				// TODO: Make it global so that version.List(...) doesn't have to be called for the same Entity if look at again since you can have multiple versions of an entity
				entityVersionList, err := version.List(destination)
				if err != nil {
					entityVersion, err = version.Latest(entityVersionList)
					if err != nil || entityVersion == "" {
						entityVersion, err = version.GeneratePseudo(destination)
					} else {
						// TODO:
						println("TODO")
					}
				}
			}*/

			/*if err != nil {
				errCh <- err
			}*/

			dirName := filepath.Base(destination)
			if strings.Contains(dirName, "@") {

			}
			//os.Rename()
			println(entityVersion)
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
