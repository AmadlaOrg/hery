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
	"gopkg.in/yaml.v3"
)

var (
	yamlUnmarshal = yaml.Unmarshal
	osReadFile    = os.ReadFile
	osStat        = os.Stat
	osIsNotExist  = os.IsNotExist
	filepathWalk  = filepath.Walk
	fileExists    = file.Exists
)

// IEntity used for mock
type IEntity interface {
	FindDir(paths storage.AbsPaths, entityVals Entity) (string, error)
	CheckDuplicate(entities []Entity, entityMeta Entity) error
	GeneratePseudoVersionPattern(name, version string) string // TODO: Move it to the version package
	CrawlDirectoriesParallel(root string) (map[string]Entity, error)
	Read(path, collectionName string) (map[string]any, error)
}

// SEntity used for mock
type SEntity struct {
	EntityVersion           version.IVersion
	EntityVersionValidation versionValidationPkg.IValidation
	EntityValidation        validation.IValidation
}

// SetSchema for appending an entity schema into the specific struct entity
// TODO: What to do if the `id` is empty
/*func (s *SEntity) setSchema(entity Entity, schema *jsonschema.Schema) error {
	var wg sync.WaitGroup
	wg.Add(len(s.Entities))

	errCh := make(chan error, 1)

	for i := range s.Entities {
		go func(i int) {
			defer wg.Done()
			if s.Entities[i].Id == "" {
				// TODO: Throw error
			}
			if s.Entities[i].Id == entity.Id {
				s.Entities[i].Schema = schema
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	return fmt.Errorf("%v", errCh) // TODO: Better error handling?
}*/

// SetContent
func (s *SEntity) setContent(entity Entity, heryContent NotFormatedContent) (Content, error) {
	// 1. Extract `_entity`
	entitySection := heryContent["_entity"].(string)
	if entity.Uri != "" {
		entitySection = entity.Uri
	} else if entitySection == "" {
		return Content{}, errors.New("no entity section found")
	}

	// 2. Extract `_id`
	// TODO: Needs to be adapted for `uuid`
	idSection := ""
	/*idSection := heryContent["_id"].(uuid)
	if entity.Id != "" {
		idSection = entity.Id
	} else if idSection == "" {
		idSection = uuid.New().String()
	}*/

	// 3. Extract `_meta`
	metaSection := heryContent["_meta"].(map[string]any)
	if metaSection == nil {
		return Content{
			Entity: entitySection,
			Id:     idSection,
		}, errors.New("_meta section is empty")
	}

	// 4. Extract `_body`
	bodySection := heryContent["_body"].(map[string]any)
	if bodySection == nil {
		return Content{
			Entity: entitySection,
			Id:     idSection,
			Meta:   metaSection,
		}, errors.New("_body section is empty")
	}

	// 5. Returns all the components of an entity content
	return Content{
		Entity: entitySection,
		Id:     idSection,
		Meta:   metaSection,
		Body:   bodySection,
	}, nil
}

// FindDir can find pseudo versioned entity directories and static versioned entities
func (s *SEntity) FindDir(paths storage.AbsPaths, entityVals Entity) (string, error) {
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
func (s *SEntity) CheckDuplicate(entities []Entity, entityMeta Entity) error {
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

// GeneratePseudoVersionPattern generates a pattern string for pseudo-versioned entities based on their name and version
// TODO: Move it?
func (s *SEntity) GeneratePseudoVersionPattern(name, version string) string {
	return fmt.Sprintf("%s@%s-*-%s", name, version[:6], version[22:])
}

// CrawlDirectoriesParallel crawls the directories in parallel and returns a map of entities
// TODO: Move it?
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

// Read makes it easy to read any yaml file with any of the two extensions: yml or yaml
// TODO: Maybe just pass collection (it might cause a cycle problem)
func (s *SEntity) Read(path, collectionName string) (map[string]any, error) {
	heryFileName := fmt.Sprintf("%s.hery", collectionName)
	heryPath := filepath.Join(path, heryFileName)

	if !fileExists(heryPath) {
		return nil, fmt.Errorf("%s does not exist", heryPath)
	}

	content, err := osReadFile(heryPath)
	if err != nil {
		return nil, err
	}

	var current map[string]interface{}
	err = yamlUnmarshal(content, &current)
	if err != nil {
		return nil, err
	}

	return current, nil
}
