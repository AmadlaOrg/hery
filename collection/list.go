package collection

import (
	"fmt"
	"github.com/AmadlaOrg/hery/storage"
	"os"
)

// List gathers all the directories that are collections in the `.henry` main storage directory to return them
func List() ([][]string, error) {
	path, err := storage.Path()
	if err != nil {
		return nil, err
	}

	// Open the directory
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	// Prepare the data for the table
	var directories [][]string
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, []string{file.Name()})
		}
	}

	return directories, nil
}
