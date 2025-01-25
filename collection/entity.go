package collection

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/LibraryUtils/url"
	entityPkg "github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/message"
	"strings"
)

type IEntityCollection interface {
	SelectByUri(entityUri string) (entityPkg.Entity, error)
	SelectById(id string) (entityPkg.Entity, error)
	Append(entity entityPkg.Entity) error
	Remove(id string) error
}

type SEntityCollection struct {
	EntityVersion           version.IVersion
	EntityVersionValidation versionValidationPkg.IValidation
	EntityValidation        validation.IValidation

	// Data
	Collection *Collection
}

// SelectAll results of an array of entities
func (s *SEntityCollection) SelectAll() *[]*entityPkg.Entity {
	return s.Collection.TransientEntities
}

// SelectByUri entity with an entity URI this functions gets the specific entity
func (s *SEntityCollection) SelectByUri(entityUri string) (entityPkg.Entity, error) {
	// 1. Set default Entity default values
	var (
		entityVals = entityPkg.Entity{
			Have:  false,
			Exist: false,
		}
		err error
	)

	s.Collection.TransientEntities = &[]*entityPkg.Entity{}

	if !s.EntityValidation.EntityUri(entityUri) {
		return entityVals, errors.New("invalid entity url")
	}

	// 2. Looks up the entity URI version
	if strings.Contains(entityUri, "@") {

		// 2.1: Extract version and if not found then it is set as the `latest`
		entityVersion, err := s.EntityVersion.Extract(entityUri)
		if err != nil {
			// 2.1.1: This error is thrown when the extractor ran into an issue
			if !errors.Is(err, version.ErrorExtractNoVersionFound) {
				return entityVals, fmt.Errorf("error extracting version: %v", err)

				// 2.1.2: Set version as `latest` if no version found
			} else {
				entityVersion = "latest"
			}
		}

		// 2.2: Removed version from the entity URI and get the full URL to the repository
		entityUriWithoutVersion := url.TrimVersion(entityUri, entityVersion)
		entityVals.RepoUrl, err = url.ExtractRepoUrl(entityUriWithoutVersion)
		if err != nil {
			return entityVals, fmt.Errorf("error extracting repo url: %v", err)
		}

		// 2.3: Extract the entity based on some entity property values
		if entityVersion == "latest" {
			// TODO: Needs goroutine
			for _, entity := range *s.Collection.TransientEntities {
				if entity.IsLatestVersion &&
					entity.RepoUrl == entityVals.RepoUrl {
					return *entity, nil
				}
			}
		} else {
			// TODO: Needs goroutine
			for _, entity := range *s.Collection.TransientEntities {
				if entity.Version == entityVersion &&
					entity.RepoUrl == entityVals.RepoUrl {
					return *entity, nil
				}
			}
		}

		// 3. If no entity URI version found
	} else {
		entityVals.RepoUrl, err = url.ExtractRepoUrl(entityUri)
		if err != nil {
			return entityVals, fmt.Errorf("error extracting repo url: %v", err)
		}

		var matchingEntities []entityPkg.Entity

		// TODO: Needs goroutine
		for _, entity := range *s.Collection.TransientEntities {
			if entity.IsLatestVersion &&
				entity.RepoUrl == entityVals.RepoUrl {
				return *entity, nil
			} else {
				matchingEntities = append(matchingEntities, *entity)
			}
		}

		matchCount := len(matchingEntities)
		if matchCount == 0 {
			return entityVals, errors.Join(
				message.ErrorNotFound,
				fmt.Errorf("no entity found with repo url %s and version %s", entityVals.RepoUrl, entityVals.Version))
		} else if matchCount >= 1 {
			return entityVals, errors.Join(
				message.ErrorMultipleFound,
				fmt.Errorf("multiple matching entities found with repo url: %s", entityVals.RepoUrl))
		}
	}

	// 4. Returns error if no entity was found
	return entityVals, errors.Join(message.ErrorNotFound, fmt.Errorf("no entity found with uri: %s", entityUri))
}

// SelectById entity
func (s *SEntityCollection) SelectById(id string) (entityPkg.Entity, error) {

	return entityPkg.Entity{}, nil
}

// Append entity to the collection
func (s *SEntityCollection) Append(entity entityPkg.Entity) error {
	// TODO:
	/*content, err := s.setContent()
	if err != nil {
		return
	}*/
	//s.Collection.TransientEntities = append(s.Collection.TransientEntities, entity)
	return nil
}

// Remove entity
func (s *SEntityCollection) Remove(id string) error {
	/*var (
		wg       sync.WaitGroup
		entities = s.Collection.TransientEntities
	)
	wg.Add(len(*entities))

	errCh := make(chan error, 1)

	for i := range *entities {
		go func(i int) {
			defer wg.Done()
			if *entities[i].Id == "" {
				errCh <- fmt.Errorf("entity id is empty")
				return
			}
			if entities[i].Id == id {
				delete(entities[i])
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	return fmt.Errorf("%v", errCh)*/

	return nil
}
