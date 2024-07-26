package get

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/build"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/AmadlaOrg/hery/util/git"
	"os"
	"sync"
)

// EntityGetter is an interface for getting entities.
type EntityGetter interface {
	Get(collectionName, storagePath string, args []string) error
	download(entityUrls []string, collectionStoragePath string) error
}

// GetterService struct implements the EntityGetter interface.
type GetterService struct {
	Git                     git.RepoManager
	EntityValidation        validation.Interface
	EntityVersion           *version.Service
	EntityVersionValidation versionValidationPkg.VersionValidator
}

// Get retrieves entities based on the provided collection name and arguments.
func (gs *GetterService) Get(storagePaths *storage.AbsPaths, entities []string) error {
	entityBuilder := build.NewEntityBuildService()
	entityBuilds := make([]entity.Entity, len(entities))
	for i, e := range entities {
		entityMeta, err := entityBuilder.MetaFromRemote(*storagePaths, e)
		if err != nil {
			return err
		}
		entityBuilds[i] = entityMeta
	}

	return gs.download(entityBuilds)
	//return nil
}

// download retrieves entities in parallel.
func (gs *GetterService) download(entitiesMeta []entity.Entity) error {
	var wg sync.WaitGroup
	wg.Add(len(entitiesMeta))

	errCh := make(chan error, len(entitiesMeta)) // Channel to collect errors
	//TODO: entityPaths := make([]string, len(entityUrls)) Maybe check with directories

	for _, entityMeta := range entitiesMeta {
		go func(entityMeta entity.Entity) {
			defer wg.Done()

			// TODO: skip (continue) loop-iteration if entity with same version was already downloaded/installed maybe make a Map

			// Create the directory if it does not exist
			err := os.MkdirAll(entityMeta.AbsPath, os.ModePerm)
			if err != nil {
				errCh <- err
			}

			// Download the Entity with `git clone`
			if err := gs.Git.FetchRepo(entityMeta.RepoUrl, entityMeta.AbsPath); err != nil {
				errCh <- fmt.Errorf("error fetching repo: %v", err)
			}
		}(entityMeta)
	}

	wg.Wait()
	close(errCh)

	var combinedErr error
	for e := range errCh {
		if combinedErr == nil {
			combinedErr = e
		} else {
			combinedErr = fmt.Errorf("%v; %v", combinedErr, e)
		}
	}

	return combinedErr
}
