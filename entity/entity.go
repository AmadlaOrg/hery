package entity

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/AmadlaOrg/hery/util/file"
	"github.com/AmadlaOrg/hery/util/url"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// IEntity
type IEntity interface {
	FindEntityDir(paths storage.AbsPaths, entityVals Entity) (string, error)
	CheckDuplicateEntity(entities []Entity, entityMeta Entity) error
	GeneratePseudoVersionPattern(name, version string) string
	CrawlDirectoriesParallel(root string) (map[string]Entity, error)
	Read(path, collectionName string) (map[string]any, error)
}

// SEntity
type SEntity struct {
	EntityVersion version.IVersion

	// Data
	Entities []Entity
}

// SetEntity
func (s *SEntity) SetEntity(entity Entity) {
	s.Entities = append(s.Entities, entity)
}

// GetEntity with an entity URI this functions gets the specific entity
func (s *SEntity) GetEntity(entityUri string) (Entity, error) {
	// 1. Set default Entity default values
	var (
		entityVals = Entity{
			Have:  false,
			Exist: false,
		}
		err error
	)

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
			for _, entity := range s.Entities {
				if entity.LatestVersion && entity.RepoUrl == entityVals.RepoUrl {
					return entity, nil
				}
			}
		} else {
			for _, entity := range s.Entities {
				if !entity.LatestVersion && entity.Version == entityVersion && entity.RepoUrl == entityVals.RepoUrl {
					return entity, nil
				}
			}
		}

		// 3. If no entity URI version found
	} else {
		entityVals.RepoUrl, err = url.ExtractRepoUrl(entityUri)
		if err != nil {
			return entityVals, fmt.Errorf("error extracting repo url: %v", err)
		}

		var matchingEntities []Entity
		for _, entity := range s.Entities {
			if entity.LatestVersion && entity.RepoUrl == entityVals.RepoUrl {
				return entity, nil
			} else {
				matchingEntities = append(matchingEntities, entity)
			}
		}

		matchCount := len(matchingEntities)
		if matchCount == 0 {
			return entityVals, errors.Join(
				ErrorNotFound,
				fmt.Errorf("no entity found with repo url %s and version %s", entityVals.RepoUrl, entityVals.Version))
		} else if matchCount >= 1 {
			return entityVals, errors.Join(
				ErrorMultipleFound,
				fmt.Errorf("multiple matching entities found with repo url: %s", entityVals.RepoUrl))
		}
	}

	// 4. Returns error if no entity was found
	return entityVals, errors.Join(ErrorNotFound, fmt.Errorf("no entity found with uri: %s", entityUri))
}

//func (s *SEntity) SetEntitySchema(entity Entity) {}

// FindEntityDir can find pseudo versioned entity directories and static versioned entities.
func (s *SEntity) FindEntityDir(paths storage.AbsPaths, entityVals Entity) (string, error) {
	entityVersionValidation := versionValidationPkg.NewEntityVersionValidationService()
	if !entityVersionValidation.PseudoFormat(entityVals.Version) {
		exactPath := entityVals.Entity

		// Check if the directory exists
		if _, err := os.Stat(exactPath); os.IsNotExist(err) {
			return "", errors.Join(
				ErrorNotFound,
				fmt.Errorf("no matching directory found for exact version: %s", exactPath))
		} else if err != nil {
			return "", err
		}

		// Return the exact path if it exists
		return exactPath, nil
	}

	// Construct the pattern
	pattern := filepath.Join(
		paths.Entities, entityVals.Origin, s.GeneratePseudoVersionPattern(entityVals.Name, entityVals.Version))

	// Use Glob to find directories matching the pattern
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", errors.Join(
			ErrorNotFound,
			fmt.Errorf("no matching directories found for pattern: %s", pattern))
	}

	if len(matches) > 1 {
		return "", errors.Join(
			ErrorMultipleFound,
			fmt.Errorf("multiple matching directories found for pattern: %s", pattern))
	}

	// Return the matched directory
	return matches[0], nil
}

// CheckDuplicateEntity checks if entityMeta is already in entityBuilds.
func (s *SEntity) CheckDuplicateEntity(entities []Entity, entityMeta Entity) error {
	entityVersionValidation := versionValidationPkg.NewEntityVersionValidationService()

	for _, existingEntity := range entities {
		if existingEntity.Origin == entityMeta.Origin &&
			existingEntity.Name == entityMeta.Name {
			if entityVersionValidation.PseudoFormat(existingEntity.Version) &&
				entityVersionValidation.PseudoFormat(entityMeta.Version) {
				// Check pseudo versions
				if s.GeneratePseudoVersionPattern(existingEntity.Name, existingEntity.Version) ==
					s.GeneratePseudoVersionPattern(entityMeta.Name, entityMeta.Version) {
					return fmt.Errorf("duplicate entity found: %v", entityMeta)
				}
			} else if existingEntity.Version == entityMeta.Version {
				// Exact version match
				return fmt.Errorf("duplicate entity found: %v", entityMeta)
			}
		}
	}
	return nil
}

// GeneratePseudoVersionPattern generates a pattern string for pseudo-versioned entities based on their name and version.
func (s *SEntity) GeneratePseudoVersionPattern(name, version string) string {
	return fmt.Sprintf("%s@%s-*-%s", name, version[:6], version[22:])
}

// CrawlDirectoriesParallel crawls the directories in parallel and returns a map of entities
func (s *SEntity) CrawlDirectoriesParallel(root string) (map[string]Entity, error) {
	entities := make(map[string]Entity)
	var mu sync.Mutex
	var wg sync.WaitGroup

	paths := make(chan string)

	// Compile regular expressions outside the loop
	versionMatchRegex := regexp.MustCompile(version.Match)
	versionFormatForDirRegex := regexp.MustCompile(version.FormatForDir)
	rePath := regexp.MustCompile(FormatEntityPathAndNameVersion)

	// Worker function
	worker := func() {
		defer wg.Done()
		for path := range paths {
			info, err := os.Stat(path)
			if err != nil {
				fmt.Println("Error stating path:", err)
				continue
			}
			if info.IsDir() {
				matched := versionMatchRegex.MatchString(info.Name())
				if matched {
					matchedPath := versionMatchRegex.MatchString(info.Name())
					if matchedPath {
						matchedPathComponents := rePath.FindStringSubmatch(path)

						// Split the directory name to extract the entity name and version
						matches := versionFormatForDirRegex.FindStringSubmatch(info.Name())
						if len(matches) == 3 && len(matchedPathComponents) == 4 {
							mu.Lock()
							entities[matches[1]] = Entity{Origin: matchedPathComponents[2], Version: matches[2]}
							mu.Unlock()
						}
					}
				}
			}
		}
	}

	// Start a fixed number of workers
	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	// Walk the directory tree
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		paths <- path
		return nil
	})
	close(paths)

	wg.Wait()
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// Read makes it easy to read any yaml file with any of the two extensions: yml or yaml
func (s *SEntity) Read(path, collectionName string) (map[string]any, error) {
	heryFileName := fmt.Sprintf("%s.hery", collectionName)
	heryPath := filepath.Join(path, heryFileName)

	if !file.Exists(heryPath) {
		return nil, fmt.Errorf("%s does not exist", heryPath)
	}

	content, err := os.ReadFile(heryPath)
	if err != nil {
		return nil, err
	}

	var current map[string]interface{}
	err = yaml.Unmarshal(content, &current)
	if err != nil {
		return nil, err
	}

	return current, nil
}
