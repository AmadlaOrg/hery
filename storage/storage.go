package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Get returns the absolute path to where the entities are stored
func Get() (string, error) {
	var entityDir string

	switch runtime.GOOS {
	case "windows":
		appDataDir := os.Getenv("APPDATA")
		entityDir = filepath.Join(appDataDir, "Hery", "entity")
	default: // "linux" and "darwin" (macOS)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error getting home directory: %s", err)
		}
		entityDir = filepath.Join(homeDir, ".hery", "entity")
	}

	return entityDir, nil
}
