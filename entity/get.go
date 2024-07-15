package entity

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/hery/collection"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	utilFilePkg "github.com/AmadlaOrg/hery/util/file"
	"github.com/AmadlaOrg/hery/util/git"
	"github.com/AmadlaOrg/hery/util/url"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

type GetInterface interface {
	Get(collectionName, storagePath string, args []string)
	download(entityUrls []string, collectionStoragePath string) error
}

type GetService struct {
	Git                     git.Interface
	EntityValidation        validation.Interface
	EntityVersion           version.Interface
	EntityVersionValidation versionValidationPkg.Interface
}

// Get with collection name and the args that are the entities urls, calls on download to get the entities
func (gs *GetService) Get(collectionName, storagePath string, args []string) {
	// Validate that all the URLs pass in the arguments are valid
	if len(args) == 0 {
		log.Fatal("No entity URL(s) specified")
	}

	for _, arg := range args {
		// TODO: Transform
		if gs.EntityValidation.EntityUrl(arg) {
			log.Fatalf("Invalid entity URL: %s", arg)
		}
	}

	collectionStoragePath := collection.Path(collectionName, storagePath)
	if !utilFilePkg.Exists(collectionStoragePath) {
		log.Fatalf("The collection storage directory does not exist: %s", collectionStoragePath)
	}

	println(collectionStoragePath)
	err := gs.download(args, collectionStoragePath)
	if err != nil {
		log.Fatal(err)
	}
}

// download in parallel all the entities
func (gs *GetService) download(entityUrls []string, collectionStoragePath string) error {
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
				/*entityVersion, err := version.Extract(entityUrl)
				if err != nil {
					errCh <- err
					return
				}*/

				// TODO: Moved to entity->validation
				var versionExists = false
				if entityVersion == "latest" {
					entityVersionList, err := gs.EntityVersion.List(entityFullRepoUrl)
					if err != nil {
						errCh <- err
						return
					}
					entityVersion, err = gs.EntityVersion.Latest(entityVersionList)
					if err != nil {
						errCh <- err
						return
					}
					versionExists = true
				} else if !gs.EntityVersionValidation.Format(entityVersion) {
					errCh <- errors.New("entity version in the entity url is wrong format")
					return
				}

				entityUrlPath = url.EntityPathUrl(entityUrl, entityVersion)
				entityFullRepoUrl = url.EntityFullRepoUrl(entityUrlPath)

				if !versionExists {
					versionExists, err := gs.EntityVersionValidation.Exists(entityUrlPath, entityVersion)
					if err != nil {
						errCh <- err
						return
					}
					if !versionExists {
						errCh <- errors.New(fmt.Sprintf("The version of the entity URL does not exist: %s", entityUrl))
						return
					}
				}
			} else {
				entityUrlPath = entityUrl
				entityFullRepoUrl = url.EntityFullRepoUrl(entityUrlPath)
				entityVersionList, err := gs.EntityVersion.List(entityFullRepoUrl)
				if err != nil {
					errCh <- err
					return
				}
				if len(entityVersionList) == 0 {
					entityVersion, err = gs.EntityVersion.GeneratePseudo(entityFullRepoUrl)
					if err != nil {
						errCh <- err
						return
					}
				} else {
					entityVersion, err = gs.EntityVersion.Latest(entityVersionList)
					if err != nil {
						errCh <- err
						return
					}
				}
				entityUrl = fmt.Sprintf("%s@%s", entityUrl, entityVersion)
			}

			destination := filepath.Join(collectionStoragePath, entityUrl)

			// Download the Entity with `git clone`
			if err := gs.Git.FetchRepo(entityFullRepoUrl, destination); err != nil {
				errCh <- fmt.Errorf("error fetching repo: %v\n", err)
			}

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
