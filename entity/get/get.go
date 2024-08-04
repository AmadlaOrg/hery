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
	"github.com/AmadlaOrg/hery/util/yaml"
	"os"
	"sync"
)

// EntityGetter is an interface for getting entities.
type EntityGetter interface {
	Get(collectionName string, storagePath string, args []string) error
	download(collectionName string, storagePaths *storage.AbsPaths, entityUrls []string, collectionStoragePath string) error
}

// GetterService struct implements the EntityGetter interface.
type GetterService struct {
	Git                     git.RepoManager
	EntityValidation        validation.Interface
	EntityVersion           *version.Service
	EntityVersionValidation versionValidationPkg.VersionValidator
	Builder                 build.MetaBuilder
}

// Get retrieves entities based on the provided collection name and arguments.
func (gs *GetterService) Get(collectionName string, storagePaths *storage.AbsPaths, entities []string) error {
	entityBuilds := make([]entity.Entity, len(entities))
	for i, e := range entities {
		entityMeta, err := gs.Builder.MetaFromRemote(*storagePaths, e)
		if err != nil {
			return err
		}

		if err = entity.CheckDuplicateEntity(entityBuilds, entityMeta); err != nil {
			return err
		}

		entityBuilds[i] = entityMeta
	}

	return gs.download(collectionName, storagePaths, entityBuilds)
}

// download retrieves entities in parallel.
func (gs *GetterService) download(collectionName string, storagePaths *storage.AbsPaths, entitiesMeta []entity.Entity) error {
	var wg sync.WaitGroup
	wg.Add(len(entitiesMeta))

	// Channel to collect errors
	errCh := make(chan error, len(entitiesMeta))

	for _, entityMeta := range entitiesMeta {
		// Skips if it is already there
		if entityMeta.Have {
			continue
		}

		go func(entityMeta entity.Entity) {
			defer wg.Done()

			// Create the directory if it does not exist
			err := os.MkdirAll(entityMeta.AbsPath, os.ModePerm)
			if err != nil {
				errCh <- err
			}

			// Download the Entity with `git clone`
			if err := gs.Git.FetchRepo(entityMeta.RepoUrl, entityMeta.AbsPath); err != nil {
				errCh <- fmt.Errorf("error fetching repo: %v", err)
			} else if !entityMeta.IsPseudoVersion {
				if err := gs.Git.CheckoutTag(entityMeta.AbsPath, entityMeta.Version); err != nil {
					errCh <- fmt.Errorf("error checking out version: %v", err)
				}
			}

			read, err := yaml.Read(entityMeta.AbsPath, collectionName)
			if err != nil {
				errCh <- fmt.Errorf("error reading yaml: %v", err)
				return
			}

			var subEntitiesMeta []entity.Entity
			for key, value := range read {
				if key == "_entity" {
					entityPath, ok := value.(string)
					if !ok {
						errCh <- fmt.Errorf("error converting yaml entity to string: %v", value)
						return
					}
					subEntityMeta, err := gs.Builder.MetaFromRemote(*storagePaths, entityPath)
					if err != nil {
						errCh <- fmt.Errorf("error fetching sub entity meta: %v", err)
						return
					}
					subEntitiesMeta = append(subEntitiesMeta, subEntityMeta)
				} else if key == "_self" {
					selfMap, ok := value.(map[string]interface{})
					if !ok {
						errCh <- fmt.Errorf("error converting yaml entity to string: %v", value)
						return
					}
					for selfKey, selfValue := range selfMap {
						if selfKey == "_entity" {
							entityPath, ok := selfValue.(string)
							if !ok {
								errCh <- fmt.Errorf("error converting yaml entity to string: %v", selfValue)
								return
							}
							subEntityMeta, err := gs.Builder.MetaFromRemote(*storagePaths, entityPath)
							if err != nil {
								errCh <- fmt.Errorf("error fetching sub entity meta: %v", err)
								return
							}
							subEntitiesMeta = append(subEntitiesMeta, subEntityMeta)
						}
					}
				}
			}

			if len(subEntitiesMeta) > 0 {
				err = gs.download(collectionName, storagePaths, subEntitiesMeta)
				if err != nil {
					errCh <- fmt.Errorf("error downloading sub entities: %v", err)
					return
				}
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
