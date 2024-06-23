package entity

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
)

// Entity holds the origin and version of an entity
type Entity struct {
	Origin  string
	Version string
}

func StorageRoot() (string, error) {
	var entityDir string

	switch runtime.GOOS {
	case "windows":
		appDataDir := os.Getenv("APPDATA")
		entityDir = filepath.Join(appDataDir, "Amadla", "entity")
	default: // "linux" and "darwin" (macOS)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error getting home directory: %s", err)
		}
		entityDir = filepath.Join(homeDir, ".amadla", "entity")
	}

	return entityDir, nil
}

func CrawlDirectories(root string) (map[string]Entity, error) {
	entities := make(map[string]Entity)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			matched, err := regexp.MatchString(`.+@v\d+\.\d+\.\d+`, info.Name())
			if err != nil {
				return err
			}
			if matched {
				// Split the directory name to extract the entity name and version
				re := regexp.MustCompile(`(.+)@v(\d+\.\d+\.\d+)`)
				matches := re.FindStringSubmatch(info.Name())
				if len(matches) == 3 {
					entities[matches[1]] = Entity{Version: matches[2]}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// CrawlDirectoriesParallel crawls the directories in parallel and returns a map of entities
func CrawlDirectoriesParallel(root string) (map[string]Entity, error) {
	entities := make(map[string]Entity)
	var mu sync.Mutex
	var wg sync.WaitGroup

	paths := make(chan string)

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
				matched, err := regexp.MatchString(`.+@v\d+\.\d+\.\d+`, info.Name())
				if err != nil {
					fmt.Println("Error matching regex:", err)
					continue
				}
				if matched {
					matchedPath, err := regexp.MatchString(`.+@v\d+\.\d+\.\d+`, info.Name())
					if err != nil {
						fmt.Println("Error matching path:", err)
					}
					if matchedPath {
						rePath := regexp.MustCompile(`^(.+[/\/]\.[A-z0-9-_]+[/\/]entity[/\/])(.+)([/\/].+@).+$`)
						matchedPathComponents := rePath.FindStringSubmatch(path)

						// Split the directory name to extract the entity name and version
						re := regexp.MustCompile(`(.+)@v(\d+\.\d+\.\d+)`)
						matches := re.FindStringSubmatch(info.Name())
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
