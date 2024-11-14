package get

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/build"
	schemaPkg "github.com/AmadlaOrg/hery/entity/schema"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/AmadlaOrg/hery/util/git"
	"os"
	"sync"
)

// TODO: Change name to Retrieve. Get is to generic and can cause confusion (e.g.: used in certain patterns)

const perm os.FileMode = os.ModePerm

// For easier mocking
var (
	osMkdirAll = os.MkdirAll
)

// IGet is an interface for getting entities.
type IGet interface {
	GetInTmp(collectionName string, entities []string) (storage.AbsPaths, error)
	Get(collectionName string, storagePaths *storage.AbsPaths, entities []string) error
	download(collectionName string, storagePaths *storage.AbsPaths, entitiesMeta []entity.Entity) error
}

// SGet struct implements the EntityGetter interface.
type SGet struct {
	Git                     git.IGit
	Entity                  entity.IEntity
	EntityValidation        validation.IValidation
	EntityVersion           version.IVersion
	EntityVersionValidation versionValidationPkg.IValidation
	Build                   build.IBuild
	Schema                  schemaPkg.ISchema
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
		entityMeta, err := s.Build.Meta(*storagePaths, e)
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

// download retrieves entities in parallel using concurrency and calls on the functions to set up, validate and
// collect sub entities
func (s *SGet) download(collectionName string, storagePaths *storage.AbsPaths, entitiesMeta []entity.Entity) error {
	var wg sync.WaitGroup
	wg.Add(len(entitiesMeta))

	// Channel to collect errors
	errCh := make(chan error, len(entitiesMeta))

	// TODO: Is Exist param ever used? Have is if it is in local and exist is for when it is found remotely

	for _, entityMeta := range entitiesMeta {
		// Skips if it is already there
		if entityMeta.Have {
			continue
		}

		go func(entityMeta entity.Entity) {
			defer wg.Done()

			// 1. Add repository in the collection directory
			err := s.addRepo(entityMeta)
			if err != nil {
				errCh <- err
				return
			}

			// 2. Gather the `<collection name>.hery` configuration file content
			heryContent, err := s.Entity.Read(entityMeta.AbsPath, collectionName)
			if err != nil {
				errCh <- fmt.Errorf("error reading yaml: %v", err)
				return
			}

			// 3. Entity validation
			// 3.1: Retrieve the `_entity`

			// 3.2: Extracted _body entity validation
			// TODO: this used to be the logic for _self
			if selfEntity := s.Schema.ExtractBody(heryContent); selfEntity != nil {

				selfEntitySchemaPath := s.Schema.GenerateSchemaPath(collectionName, entityMeta.AbsPath)
				selfEntitySchema, err := s.Schema.Load(selfEntitySchemaPath)
				if err != nil {
					errCh <- fmt.Errorf("error loading schema: %v", err)
					return
				}

				err = s.EntityValidation.Entity(collectionName, selfEntitySchema, selfEntity)
				if err != nil {
					errCh <- fmt.Errorf("error validating entity: %v", err)
					return
				}

				delete(heryContent, "_body")

				err = s.EntityValidation.Entity(collectionName, selfEntitySchema, selfEntity)
				if err != nil {
					errCh <- fmt.Errorf("error validating entity: %v", err)
					return
				}
			} else {
				// No `_body`
				// TODO: Double check

			}

			// TODO: Add validation to the main entity

			// 3. Validate the content of hery file content to make sure it does not cause is issue later in the code
			//
			// -- This follows the Fail Fast principal --
			//
			// TODO:
			/*err = s.EntityValidation.Entity(entityMeta.AbsPath, collectionName, entityMeta.Entity, heryContent)
			if err != nil {
				errCh <- fmt.Errorf("error validating entity: %v", err)
				return
			}*/

			// 4. The reference to the other entities are found in the hery file content
			//
			// This function gathers the `_entity` properties that have the entity URIs that are used to pull the entity
			// repositories.
			//
			// TODO: Add limit on how many times this can be called since we don't want infinite loop (maybe add a counter)
			err = s.collectSubEntities(collectionName, storagePaths, heryContent)
			if err != nil {
				errCh <- fmt.Errorf("error collecting sub entities: %v", err)
				return
			}

		}(entityMeta)
	}

	wg.Wait() // TODO: Just hangs here
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

// addRepo does all the tasks required to setup a new entity (or entity with a different version)
// TODO: Add a timer limit (maybe go-git has something for that) so that it does not get
// TODO: Make sure we have clear error. Because it seems it just hangs without clear error.
// TODO: Might want to add hashing of the entity once it was downloaded to have verification that nothing was corrupted for Fail Fast principal
func (s *SGet) addRepo(entityMeta entity.Entity) error {
	// 1. Create the directory if it does not exist
	err := osMkdirAll(entityMeta.AbsPath, perm)
	if err != nil {
		return err
	}

	// 2. Download the Entity with `git clone`
	if err = s.Git.FetchRepo(entityMeta.RepoUrl, entityMeta.AbsPath); err != nil {
		return fmt.Errorf("error fetching repo: %v", err)
	}

	// 3. Changes the repository to the tag (version) that was pass
	if !entityMeta.IsPseudoVersion {
		if err = s.Git.CheckoutTag(entityMeta.AbsPath, entityMeta.Version); err != nil {
			return fmt.Errorf("error checking out version: %v", err)
		}
	}

	return nil
}

// collectSubEntities Calls on download function with the entity URIs that were found in the `_entity`
//
// For the `_body` contains the initial configuration of the setup of the entity `.hery` configuration. It might contain
// `_entity` and this function also pulls those sub entities.
//
// download function is call because inside any entities there might be again sub entities.
// TODO: Needs to be reviewed based on the new structure
func (s *SGet) collectSubEntities(
	collectionName string,
	storagePaths *storage.AbsPaths,
	henryContent map[string]interface{}) error {

	// 1. Loops through the properties found in the `.hery` configuration file
	// found in the `_entity` or the `_entity` in `_body`
	var subEntitiesMeta []entity.Entity
	for key, value := range henryContent {
		// TODO: There is multiple keys that are at the same level that can be present: _entity, _meta, _id, and _body
		if key == "_entity" {
			entityPath, ok := value.(string)
			if !ok {
				return fmt.Errorf("error converting yaml entity to string: %v", value)
			}
			subEntityMeta, err := s.Build.Meta(*storagePaths, entityPath)
			if err != nil {
				return fmt.Errorf("error fetching sub entity meta: %v", err)
			}
			subEntitiesMeta = append(subEntitiesMeta, subEntityMeta)
		} else if key == "_body" {
			selfMap, ok := value.(map[string]interface{})
			if !ok {
				return fmt.Errorf("error converting yaml entity to string: %v", value)
			}
			for selfKey, selfValue := range selfMap {
				if selfKey == "_entity" {
					entityPath, ok := selfValue.(string)
					if !ok {
						return fmt.Errorf("error converting yaml entity to string: %v", selfValue)
					}
					subEntityMeta, err := s.Build.Meta(*storagePaths, entityPath)
					if err != nil {
						return fmt.Errorf("error fetching sub entity meta: %v", err)
					}
					subEntitiesMeta = append(subEntitiesMeta, subEntityMeta)
				}
			}
		}
	}

	// 2. If sub entities found then send it to download function
	// TODO: Add check that this entity does not already in `Have == true`
	// TODO: Or maybe add that logic at a higher level so that it is not added to the `subEntitiesMeta` list
	if len(subEntitiesMeta) > 0 {
		err := s.download(collectionName, storagePaths, subEntitiesMeta)
		if err != nil {
			return fmt.Errorf("error downloading sub entities: %v", err)
		}
	}

	return nil
}
