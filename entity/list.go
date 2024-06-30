package entity

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity/version"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

// CrawlDirectories
func CrawlDirectories(root string) (map[string]Entity, error) {
	entities := make(map[string]Entity)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			matched, err := regexp.MatchString(version.Match, info.Name())
			if err != nil {
				return err
			}
			if matched {
				// Split the directory name to extract the entity name and version
				re := regexp.MustCompile(version.FormatForDir)
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

	// Compile regular expressions outside the loop
	versionMatchRegex := regexp.MustCompile(version.Match)
	versionFormatForDirRegex := regexp.MustCompile(version.FormatForDir)
	rePath := regexp.MustCompile(formatEntityPathAndNameVersion)

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
