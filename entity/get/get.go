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

// IGet is an interface for getting entities.
type IGet interface {
	GetInTmp(collectionName string, entities []string) (storage.AbsPaths, error)
	Get(collectionName string, storagePath string, args []string) error
	download(collectionName string, storagePaths *storage.AbsPaths, entitiesMeta []entity.Entity) error
}

// SGet struct implements the EntityGetter interface.
type SGet struct {
	Git                     git.RepoManager
	Entity                  entity.IEntity
	EntityValidation        validation.IValidation
	EntityVersion           version.IVersion
	EntityVersionValidation versionValidationPkg.IValidation
	Build                   build.IBuild
}

// GetInTmp retrieves entities based on the provided collection name and entities
func (s *SGet) GetInTmp(collectionName string, entities []string) (storage.AbsPaths, error) {
	storageService := storage.NewStorageService()

	// Replace paths with temporary directory before .<collectionName>
	storagePaths, err := storageService.TmpPaths(collectionName)
	if err != nil {
		return storage.AbsPaths{}, err
	}

	err = storageService.MakePaths(*storagePaths)
	if err != nil {
		return *storagePaths, err
	}

	err = s.Get(collectionName, storagePaths, entities)
	if err != nil {
		return *storagePaths, err
	}

	return *storagePaths, nil
}

// Get retrieves entities based on the provided collection name and entities
func (s *SGet) Get(collectionName string, storagePaths *storage.AbsPaths, entities []string) error {
	entityBuilds := make([]entity.Entity, len(entities))
	for i, e := range entities {
		entityMeta, err := s.Build.MetaFromRemote(*storagePaths, e)
		if err != nil {
			return err
		}

		if err = s.Entity.CheckDuplicateEntity(entityBuilds, entityMeta); err != nil {
			return err
		}

		entityBuilds[i] = entityMeta
	}

	return s.download(collectionName, storagePaths, entityBuilds)
}

// download retrieves entities in parallel.
func (s *SGet) download(collectionName string, storagePaths *storage.AbsPaths, entitiesMeta []entity.Entity) error {
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
				return
			}

			// Download the Entity with `git clone`
			if err := s.Git.FetchRepo(entityMeta.RepoUrl, entityMeta.AbsPath); err != nil {
				errCh <- fmt.Errorf("error fetching repo: %v", err)
				return
			}

			// Changes the repository to the tag (version) that was pass
			if !entityMeta.IsPseudoVersion {
				if err := s.Git.CheckoutTag(entityMeta.AbsPath, entityMeta.Version); err != nil {
					errCh <- fmt.Errorf("error checking out version: %v", err)
					return
				}
			}

			read, err := yaml.Read(entityMeta.AbsPath, collectionName)
			if err != nil {
				errCh <- fmt.Errorf("error reading yaml: %v", err)
				return
			}

			err = s.EntityValidation.Entity(collectionName, entityMeta.AbsPath)
			if err != nil {
				errCh <- fmt.Errorf("error validating entity: %v", err)
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
					subEntityMeta, err := s.Build.MetaFromRemote(*storagePaths, entityPath)
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
							subEntityMeta, err := s.Build.MetaFromRemote(*storagePaths, entityPath)
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
				err = s.download(collectionName, storagePaths, subEntitiesMeta)
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
