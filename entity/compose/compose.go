package compose

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/heryext"
	"github.com/AmadlaOrg/hery/storage"
	utilObjectPkg "github.com/AmadlaOrg/hery/util/object"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

// IComposer is an interface for composing entities.
type IComposer interface {
	ComposeEntity(entityArg string, printToScreen bool) error
	parseEntityArg(entityArg string) (string, string, error)
	findEntityDirParallel(root, name, version string) (string, error)
	mergeYamlFiles(dir string) ([]byte, error)
}

// SComposer struct implements the EntityComposer interface.
type SComposer struct {
	Storage storage.IStorage
	HeryExt heryext.IHeryExt
}

// ComposeEntity gathers as many details about an Entity as possible and composes it.
func (s *SComposer) ComposeEntity(entityArg string, printToScreen bool) error {
	root, err := s.Storage.Main()
	if err != nil {
		return err
	}

	entityName, version, err := s.parseEntityArg(entityArg)
	if err != nil {
		return err
	}

	/*entities, err := CrawlDirectoriesParallel(entityDir)
	if err != nil {
		return err
	}*/

	// Find the directory
	entityDir, err := s.findEntityDirParallel(root, entityName, version)
	if err != nil {
		return err
	}

	fmt.Println("entityDir: " + entityDir)

	// Read and merge the YAML files
	mergedYaml, err := s.mergeYamlFiles(entityDir)
	if err != nil {
		return err
	}

	if printToScreen {
		fmt.Println(string(mergedYaml))
	} else {
		err = os.WriteFile(entityName+".lock", mergedYaml, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SComposer) parseEntityArg(entityArg string) (string, string, error) {
	entityNamePattern := regexp.MustCompile(entity.EntityNameMatch)
	entityWithVersionPattern := regexp.MustCompile(entity.EntityNameAndVersionMatch)

	if entityWithVersionPattern.MatchString(entityArg) {
		matches := entityWithVersionPattern.FindStringSubmatch(entityArg)
		return matches[1], matches[2], nil
	} else if entityNamePattern.MatchString(entityArg) {
		return entityArg, "", nil
	}

	return "", "", fmt.Errorf("invalid entity argument: %s", entityArg)
}

func (s *SComposer) findEntityDirParallel(root, name, version string) (string, error) {
	var matchedDir string
	var mu sync.Mutex
	var wg sync.WaitGroup
	var once sync.Once
	done := make(chan struct{})

	paths := make(chan string, 100) // Buffered channel to avoid blocking

	// Worker function
	worker := func() {
		defer wg.Done()
		for {
			select {
			case path, ok := <-paths:
				if !ok {
					return
				}
				info, err := os.Stat(path)
				if err != nil {
					continue
				}
				if info.IsDir() {
					matched, err := regexp.MatchString(fmt.Sprintf(`%s@v\d+\.\d+\.\d+`, name), info.Name())
					if err != nil {
						continue
					}
					if matched {
						if version == "" || info.Name() == name+version {
							mu.Lock()
							matchedDir = path
							readYaml, err := s.HeryExt.Read(path, "amadla")
							if err != nil {
								mu.Unlock()
								return
							}
							marshalled, err := yaml.Marshal(readYaml)
							if err != nil {
								mu.Unlock()
								return
							}
							fmt.Printf("%s\n", marshalled)
							mu.Unlock()
							once.Do(func() { close(done) })
							return
						}
					}
				}
			case <-done:
				return
			}
		}
	}

	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		select {
		case paths <- path:
		case <-done:
			return filepath.SkipDir
		}
		return nil
	})
	close(paths)
	wg.Wait()

	if err != nil {
		return "", err
	}
	if matchedDir == "" {
		return "", fmt.Errorf("entity %s%s not found", name, version)
	}
	return matchedDir, nil
}

func (s *SComposer) mergeYamlFiles(dir string) ([]byte, error) {
	var merged map[string]interface{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".yaml") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var current map[string]interface{}
			err = yaml.Unmarshal(content, &current)
			if err != nil {
				return err
			}
			if merged == nil {
				merged = current
			} else {
				merged = utilObjectPkg.MergeMultilevel(merged, current, true)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(merged)
}
