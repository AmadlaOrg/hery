package entity

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/AmadlaOrg/hery/util/file"
	"github.com/AmadlaOrg/hery/util/url"
	"github.com/google/uuid"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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
	Add(entity Entity)
	Get(entityUri string) (Entity, error)
	GetAll() []Entity
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

	// Data
	Entities []Entity
}

// Add for appending an entity into the struct entity list
func (s *SEntity) Add(entity Entity) error {
	content, err := s.setContent()
	if err != nil {
		return
	}

	s.Entities = append(s.Entities, entity)
}

// Get with an entity URI this functions gets the specific entity
func (s *SEntity) Get(entityUri string) (Entity, error) {

	// 1. Set default Entity default values
	var (
		entityVals = Entity{
			Have:  false,
			Exist: false,
		}
		err error
	)

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
			for _, entity := range s.Entities {
				if entity.LatestVersion &&
					entity.RepoUrl == entityVals.RepoUrl {
					return entity, nil
				}
			}
		} else {
			// TODO: Needs goroutine
			for _, entity := range s.Entities {
				if entity.Version == entityVersion &&
					entity.RepoUrl == entityVals.RepoUrl {
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

		// TODO: Needs goroutine
		for _, entity := range s.Entities {
			if entity.LatestVersion &&
				entity.RepoUrl == entityVals.RepoUrl {
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

// GetAll results of an array of entities
func (s *SEntity) GetAll() []Entity {
	return s.Entities
}

// SetSchema for appending an entity schema into the specific struct entity
// TODO: What to do if the `id` is empty
func (s *SEntity) setSchema(entity Entity, schema *jsonschema.Schema) error {
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
}

// SetContent
func (s *SEntity) setContent(entity Entity, heryContent NotFormatedContent) (Content, error) {
	// 1. Extract `_entity`
	entitySection := heryContent["_entity"].(string)
	if entity.Entity != "" {
		entitySection = entity.Entity
	} else if entitySection == "" {
		return Content{}, errors.New("no entity section found")
	}

	// 2. Extract `_id`
	idSection := heryContent["_id"].(string)
	if entity.Id != "" {
		idSection = entity.Id
	} else if idSection == "" {
		idSection = uuid.New().String()
	}

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

// FindDir can find pseudo versioned entity directories and static versioned entities.
func (s *SEntity) FindDir(paths storage.AbsPaths, entityVals Entity) (string, error) {
	if !s.EntityVersionValidation.PseudoFormat(entityVals.Version) {
		exactPath := entityVals.Entity

		// Check if the directory exists
		if _, err := osStat(exactPath); osIsNotExist(err) {
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
