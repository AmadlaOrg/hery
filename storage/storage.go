package storage

import (
	"fmt"
	"github.com/AmadlaOrg/hery/util/file"
	"os"
	"path/filepath"
	"runtime"
)

// Path returns the absolute path to where the entities are stored
func Path() (string, error) {
	var entityDir string

	//
	// Using env var
	//

	envStoragePathValue := os.Getenv(heryStoragePath)

	if envStoragePathValue != "" {
		envStoragePath, err := filepath.Abs(envStoragePathValue)
		if err != nil {
			return "", err
		}
		return envStoragePath, nil
	}

	//
	// Using current location
	//

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	localStoragePath := filepath.Join(cwd, ".hery")

	if file.Exists(localStoragePath) {
		return localStoragePath, nil
	}

	//
	// Default
	//

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
