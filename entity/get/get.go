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

// IGet is an interface for getting entities.
type IGet interface {
	GetInTmp(entities []string) (storage.AbsPaths, error)
	Get(storagePaths *storage.AbsPaths, entities []string) error
	download(storagePaths *storage.AbsPaths, entitiesMeta []entity.Entity) error
}

// SGet struct implements the EntityGetter interface.
type SGet struct {
	Entity                  entity.IEntity
	EntityValidation        validation.IValidation
	EntityVersion           version.IVersion
	EntityVersionValidation versionValidationPkg.IValidation
	Build                   build.IBuild
	Schema                  schemaPkg.ISchema

	// Config
	GitConfig *gitConfig.Config
}

const perm os.FileMode = os.ModePerm

// For easier mocking
var (
	osMkdirAll       = os.MkdirAll
	gitNewGitService = git.NewGitService
)

// GetInTmp retrieves entities into a temporary directory
func (s *SGet) GetInTmp(entities []string) (storage.AbsPaths, error) {
	storageService := storage.NewStorageService()

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
func (s *SGet) Get(storagePaths *storage.AbsPaths, entities []string) error {
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
func (s *SGet) download(storagePaths *storage.AbsPaths, entitiesMeta []entity.Entity) error {
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

			// 2b. Resolve _parent chains via deep merge
			parentLookup := func(selfURI string) (map[string]any, error) {
				parentMeta, lookupErr := s.Build.Meta(*storagePaths, selfURI)
				if lookupErr != nil {
					return nil, lookupErr
				}
				if !parentMeta.Have {
					if dlErr := s.download(storagePaths, []entity.Entity{parentMeta}); dlErr != nil {
						return nil, dlErr
					}
				}
				parentDocs, readErr := s.Entity.ReadAll(parentMeta.AbsPath)
				if readErr != nil || len(parentDocs) == 0 {
					return nil, fmt.Errorf("could not read parent entity %s", selfURI)
				}
				return parentDocs[0], nil
			}

			for i, doc := range documents {
				resolved, resolveErr := merge.ResolveParentChain(doc, parentLookup)
				if resolveErr != nil {
					errCh <- fmt.Errorf("error resolving _parent chain: %v", resolveErr)
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
func (s *SGet) addRepo(entityMeta entity.Entity) error {
	err := osMkdirAll(entityMeta.AbsPath, perm)
	if err != nil {
		return err
	}

	gitService := gitNewGitService(entityMeta.RepoUrl, entityMeta.AbsPath, s.GitConfig)

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
