package storage

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_Integration_Main(t *testing.T) {
	storageService := NewStorageService()
	paths, err := storageService.Paths()
	if err != nil {
		t.Fatal("Failed to get paths")
	}

	// Determine the main directory based on the OS
	var mainDir string
	switch runtime.GOOS {
	case "windows":
		appDataDir := os.Getenv("LOCALAPPDATA")
		if appDataDir == "" {
			appDataDir = os.Getenv("APPDATA")
		}
		mainDir = filepath.Join(appDataDir, "hery")
	default:
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			homeDir, _ := os.UserHomeDir()
			cacheDir = filepath.Join(homeDir, ".cache")
		}
		mainDir = filepath.Join(cacheDir, "hery")
	}

	expectedPaths := AbsPaths{
		Storage:  mainDir,
		Entities: filepath.Join(mainDir, "entity"),
	}

	if paths.Storage != expectedPaths.Storage {
		t.Errorf("Storage path mismatch. Got: %s, Want: %s", paths.Storage, expectedPaths.Storage)
	}
	if paths.Entities != expectedPaths.Entities {
		t.Errorf("Entities path mismatch. Got: %s, Want: %s", paths.Entities, expectedPaths.Entities)
	}
}
