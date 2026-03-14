package get

import (
	"fmt"
	"github.com/AmadlaOrg/LibraryUtils/git"
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/build"
	"github.com/AmadlaOrg/hery/entity/merge"
	schemaPkg "github.com/AmadlaOrg/hery/entity/schema"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/storage"
	"os"
	"sync"
)

// Getter is an interface for getting entities.
type Getter interface {
	GetInTmp(entities []string) (storage.AbsPaths, error)
	Get(storagePaths *storage.AbsPaths, entities []string) error
	download(storagePaths *storage.AbsPaths, entitiesMeta []entity.Entity) error
}

// getter struct implements the EntityGetter interface.
type getter struct {
	Entity                  entity.Service
	EntityValidation        validation.Validator
	EntityVersion           version.Version
	EntityVersionValidation versionValidationPkg.Validator
	Build                   build.Builder
	Schema                  schemaPkg.Schema

	// Config
	GitConfig *gitConfig.Config
}

const perm os.FileMode = os.ModePerm

// For easier mocking
var (
	osMkdirAll       = os.MkdirAll
	gitNew = git.New
)

// GetInTmp retrieves entities into a temporary directory
func (s *getter) GetInTmp(entities []string) (storage.AbsPaths, error) {
	storageService := storage.New()

	storagePaths, err := storageService.TmpPaths()
	if err != nil {
		return storage.AbsPaths{}, err
	}

	err = storageService.MakePaths(*storagePaths)
	if err != nil {
		return *storagePaths, err
	}

	err = s.Get(storagePaths, entities)
	if err != nil {
		return *storagePaths, err
	}

	return *storagePaths, nil
}

// Get retrieves entities based on the provided entity URIs
func (s *getter) Get(storagePaths *storage.AbsPaths, entities []string) error {
	entityBuilds := make([]entity.Entity, len(entities))
	for i, e := range entities {
		entityMeta, err := s.Build.Meta(*storagePaths, e)
		if err != nil {
			return err
		}

		if err = s.Entity.CheckDuplicate(entityBuilds, entityMeta); err != nil {
			return err
		}

		entityBuilds[i] = entityMeta
	}

	return s.download(storagePaths, entityBuilds)
}

// download retrieves entities in parallel
func (s *getter) download(storagePaths *storage.AbsPaths, entitiesMeta []entity.Entity) error {
	var wg sync.WaitGroup
	wg.Add(len(entitiesMeta))

	errCh := make(chan error, 1)

	for _, entityMeta := range entitiesMeta {
		if entityMeta.Have {
			wg.Done()
			continue
		}

		go func(entityMeta entity.Entity) {
			defer wg.Done()

			// 1. Clone the repository
			err := s.addRepo(entityMeta)
			if err != nil {
				errCh <- err
				return
			}

			// 2. Read all .hery files in the entity directory
			documents, err := s.Entity.ReadAll(entityMeta.AbsPath)
			if err != nil {
				errCh <- fmt.Errorf("error reading hery files: %v", err)
				return
			}

			// 2b. Resolve _extends chains via deep merge
			extendsLookup := func(selfURI string) (map[string]any, error) {
				extendsMeta, lookupErr := s.Build.Meta(*storagePaths, selfURI)
				if lookupErr != nil {
					return nil, lookupErr
				}
				if !extendsMeta.Have {
					if dlErr := s.download(storagePaths, []entity.Entity{extendsMeta}); dlErr != nil {
						return nil, dlErr
					}
				}
				extendsDocs, readErr := s.Entity.ReadAll(extendsMeta.AbsPath)
				if readErr != nil || len(extendsDocs) == 0 {
					return nil, fmt.Errorf("could not read extended entity %s", selfURI)
				}
				return extendsDocs[0], nil
			}

			for i, doc := range documents {
				resolved, resolveErr := merge.ResolveExtendsChain(doc, extendsLookup)
				if resolveErr != nil {
					errCh <- fmt.Errorf("error resolving _extends chain: %v", resolveErr)
					return
				}
				documents[i] = resolved
			}

			// 3. Collect sub-entities referenced via _type in _body (via schema $ref)
			for _, doc := range documents {
				if typeVal, ok := doc["_type"].(string); ok {
					subEntityMeta, err := s.Build.Meta(*storagePaths, typeVal)
					if err != nil {
						errCh <- fmt.Errorf("error fetching sub entity meta: %v", err)
						return
					}
					if !subEntityMeta.Have {
						if err := s.download(storagePaths, []entity.Entity{subEntityMeta}); err != nil {
							errCh <- fmt.Errorf("error downloading sub entities: %v", err)
							return
						}
					}
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

// addRepo clones the entity repository and checks out the correct version
func (s *getter) addRepo(entityMeta entity.Entity) error {
	err := osMkdirAll(entityMeta.AbsPath, perm)
	if err != nil {
		return err
	}

	gitService := gitNew(entityMeta.RepoUrl, entityMeta.AbsPath, s.GitConfig)

	if err = gitService.Clone(); err != nil {
		return fmt.Errorf("error fetching repo: %v", err)
	}

	if !entityMeta.IsPseudoVersion {
		if err = gitService.CheckoutTag(entityMeta.Version); err != nil {
			return fmt.Errorf("error checking out version: %v", err)
		}
	}

	return nil
}
