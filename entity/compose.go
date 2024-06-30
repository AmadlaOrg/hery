package entity

import (
	"fmt"
	"github.com/AmadlaOrg/hery/storage"
	utilObjectPkg "github.com/AmadlaOrg/hery/util/object"
	utilYamlPkg "github.com/AmadlaOrg/hery/util/yaml"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

func ComposeEntity(entityArg string, printToScreen bool) error {
	root, err := storage.Path()
	if err != nil {
		return err
	}

	entityName, version, err := parseEntityArg(entityArg)
	if err != nil {
		return err
	}

	/*entities, err := CrawlDirectoriesParallel(entityDir)
	if err != nil {
		return err
	}*/

	// Find the directory
	entityDir, err := findEntityDirParallel(root, entityName, version)
	if err != nil {
		return err
	}

	println("entityDir: " + entityDir)

	// Read and merge the YAML files
	/*mergedYaml, err := mergeYamlFiles(entityDir)
	if err != nil {
		return err
	}

	if printToScreen {
		fmt.Println(string(mergedYaml))
	} else {
		err = ioutil.WriteFile(entityName+".lock", mergedYaml, 0644)
		if err != nil {
			return err
		}
	}*/

	return nil
}

func parseEntityArg(entityArg string) (string, string, error) {
	// Validate entity name and version separately
	entityNamePattern := regexp.MustCompile(entityNameMatch)
	entityWithVersionPattern := regexp.MustCompile(entityNameAndVersionMatch)

	if entityWithVersionPattern.MatchString(entityArg) {
		matches := entityWithVersionPattern.FindStringSubmatch(entityArg)
		return matches[1], matches[2], nil
	} else if entityNamePattern.MatchString(entityArg) {
		return entityArg, "", nil
	}

	return "", "", fmt.Errorf("invalid entity argument: %s", entityArg)
}

func findEntityDirParallel(root, name, version string) (string, error) {
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
					fmt.Println("Error stating path:", err)
					continue
				}
				if info.IsDir() {
					matched, err := regexp.MatchString(fmt.Sprintf(`%s@v\d+\.\d+\.\d+`, name), info.Name())
					if err != nil {
						fmt.Println("Error matching regex:", err)
						continue
					}
					if matched {
						if version == "" || info.Name() == name+version {
							mu.Lock()

							matchedDir = path
							readYaml, err := utilYamlPkg.Read(path, "amadla")
							if err != nil {
								return
							}
							marshalled, err := yaml.Marshal(readYaml)
							if err != nil {
								fmt.Println("Error marshalling YAML:", err)
								return
							}
							fmt.Printf("%s\n", marshalled)

							mu.Unlock()
							once.Do(func() { close(done) }) // Signal to stop processing
							return
						}
					}
				}
			case <-done:
				return
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
		select {
		case paths <- path:
		case <-done: // Exit early if we already found a match
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

func mergeYamlFiles(dir string) ([]byte, error) {
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
