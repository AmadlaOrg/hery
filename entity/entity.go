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

type IEntity interface {
	FindEntityDir(paths storage.AbsPaths, entityVals Entity) (string, error)
	CheckDuplicateEntity(entities []Entity, entityMeta Entity) error
	GeneratePseudoVersionPattern(name, version string) string
	CrawlDirectoriesParallel(root string) (map[string]Entity, error)
	Read(path, collectionName string) (map[string]any, error)
}

type SEntity struct {
	EntityVersion version.IVersion

	// Data
	Entities []Entity
}

// GetEntity
func (s *SEntity) GetEntity(entityUri string) (Entity, error) {
	var (
		entityVals = Entity{
			Have:  false,
			Exist: false,
		}
		//err error
	)

	if strings.Contains(entityUri, "@") {
		entityVersion, err := s.EntityVersion.Extract(entityUri)
		if err != nil {
			if !errors.Is(err, version.ErrorExtractNoVersionFound) {
				return entityVals, fmt.Errorf("error extracting version: %v", err)
			} else {
				entityVersion = "latest"
			}
		}

		entityUriWithoutVersion := url.TrimVersion(entityUri, entityVersion)
		entityVals.RepoUrl, err = url.ExtractRepoUrl(entityUriWithoutVersion)
		if err != nil {
			return entityVals, fmt.Errorf("error extracting repo url: %v", err)
		}

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
	} else {
		var (
			entityVals = entity.Entity{
				Have:  false,
				Exist: false,
			}
			err error
		)

		entityUriWithoutVersion := url.TrimVersion(entityUri, entityVersion)
		entityVals.RepoUrl, err = url.ExtractRepoUrl(entityUriWithoutVersion)

		if err != nil {
			return entityVals, fmt.Errorf("error extracting repo url: %v", err)
		}

		entityVersionList, err := s.EntityVersion.List(entityVals.RepoUrl)
		if err != nil {
			return entityVals, fmt.Errorf("error listing versions: %v", err)
		}
	}

	/**/
	return Entity{}, fmt.Errorf("no entity found with uri: %s", entityUri)
}

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
