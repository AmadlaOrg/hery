package entity

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/AmadlaOrg/LibraryUtils/file"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/message"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/goccy/go-yaml"
)

var (
	yamlUnmarshal = yaml.Unmarshal
	osReadFile    = os.ReadFile
	osStat        = os.Stat
	osIsNotExist  = os.IsNotExist
	filepathWalk  = filepath.Walk
	fileExists    = file.Exists
)

// Service used for mock
type Service interface {
	FindDir(paths storage.AbsPaths, entityVals Entity) (string, error)
	CheckDuplicate(entities []Entity, entityMeta Entity) error
	GeneratePseudoVersionPattern(name, version string) string
	CrawlDirectoriesParallel(root string) (map[string]Entity, error)
	ReadAll(path string) ([]map[string]any, error)
}

// service used for mock
type service struct {
	EntityVersion           version.Version
	EntityVersionValidation versionValidationPkg.Validator
	EntityValidation        validation.Validator
}

// setContent extracts the five reserved HERY properties from raw YAML content.
func (s *service) setContent(entity Entity, heryContent NotFormatedContent) (Content, error) {
	// 1. Extract `_type` (required)
	typeSection, _ := heryContent["_type"].(string)
	if entity.Uri != "" {
		typeSection = entity.Uri
	}
	if typeSection == "" {
		return Content{}, errors.New("_type is required")
	}

	// 2. Extract `_extends` (optional)
	extendsSection, _ := heryContent["_extends"].(string)

	// 3. Extract `_meta` (optional)
	metaSection, _ := heryContent["_meta"].(map[string]any)

	// 4. Extract `_body` (optional)
	bodySection, _ := heryContent["_body"].(map[string]any)

	// 5. Extract `_requires` (optional)
	var requiresSection []string
	if rawRequires, ok := heryContent["_requires"].([]any); ok {
		for _, item := range rawRequires {
			if s, ok := item.(string); ok {
				requiresSection = append(requiresSection, s)
			}
		}
	}

	return Content{
		Type:     typeSection,
		Extends:  extendsSection,
		Meta:     metaSection,
		Body:     bodySection,
		Requires: requiresSection,
	}, nil
}

// FindDir can find pseudo versioned entity directories and static versioned entities
func (s *service) FindDir(paths storage.AbsPaths, entityVals Entity) (string, error) {
	if !s.EntityVersionValidation.PseudoFormat(entityVals.Version) {
		exactPath := entityVals.Uri

		// Check if the directory exists
		if _, err := osStat(exactPath); osIsNotExist(err) {
			return "", errors.Join(
				message.ErrorNotFound,
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
			message.ErrorNotFound,
			fmt.Errorf("no matching directories found for pattern: %s", pattern))
	}

	if len(matches) > 1 {
		return "", errors.Join(
			message.ErrorMultipleFound,
			fmt.Errorf("multiple matching directories found for pattern: %s", pattern))
	}

	// Return the matched directory
	return matches[0], nil
}

// CheckDuplicate checks if entityMeta is already in entityBuilds.
func (s *service) CheckDuplicate(entities []Entity, entityMeta Entity) error {
	for _, existingEntity := range entities {
		if existingEntity.Origin == entityMeta.Origin &&
			existingEntity.Name == entityMeta.Name {
			if s.EntityVersionValidation.PseudoFormat(existingEntity.Version) &&
				s.EntityVersionValidation.PseudoFormat(entityMeta.Version) {
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

// GeneratePseudoVersionPattern delegates to version.GeneratePseudoPattern.
func (s *service) GeneratePseudoVersionPattern(name, ver string) string {
	return s.EntityVersion.GeneratePseudoPattern(name, ver)
}

// CrawlDirectoriesParallel crawls the directories in parallel and returns a map of entities.
func (s *service) CrawlDirectoriesParallel(root string) (map[string]Entity, error) {
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
			info, err := osStat(path)
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
	err := filepathWalk(root, func(path string, info os.FileInfo, err error) error {
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

// ReadAll reads all .hery files in a directory and returns each YAML document as a map.
// Supports multi-document YAML files (separated by ---).
func (s *service) ReadAll(dir string) ([]map[string]any, error) {
	var documents []map[string]any

	err := filepathWalk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || filepath.Ext(path) != ".hery" {
			return nil
		}

		content, readErr := osReadFile(path)
		if readErr != nil {
			return readErr
		}

		var doc map[string]any
		if unmarshalErr := yamlUnmarshal(content, &doc); unmarshalErr != nil {
			return unmarshalErr
		}
		if doc != nil {
			documents = append(documents, doc)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return documents, nil
}
