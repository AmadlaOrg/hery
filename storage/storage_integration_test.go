package storage

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_Integration_Main(t *testing.T) {
	storageService := NewStorageService()
	paths, err := storageService.Paths("test")
	if err != nil {
		t.Fatal("Failed to get paths")
	}

	// Determine the main directory based on the OS
	var mainDir string
	switch runtime.GOOS {
	case "windows":
		appDataDir := os.Getenv("APPDATA")
		mainDir = filepath.Join(appDataDir, "Hery")
	default: // "linux" and "darwin" (macOS)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Fatal("Failed to get user home dir")
		}
		mainDir = filepath.Join(homeDir, ".hery")
	}

	// Define expected paths based on the main directory
	expectedPaths := AbsPaths{
		Storage:    mainDir,
		Catalog:    filepath.Join(mainDir, "collection"),
		Collection: filepath.Join(mainDir, "collection", "test"),
		Entities:   filepath.Join(mainDir, "collection", "test", "entity"),
		Cache:      filepath.Join(mainDir, "collection", "test", "test.cache"),
	}

	// Check if the actual paths match the expected paths
	if paths.Storage != expectedPaths.Storage {
		t.Errorf("Storage path mismatch. Got: %s, Want: %s", paths.Storage, expectedPaths.Storage)
	}
	if paths.Catalog != expectedPaths.Catalog {
		t.Errorf("Catalog path mismatch. Got: %s, Want: %s", paths.Catalog, expectedPaths.Catalog)
	}
	if paths.Collection != expectedPaths.Collection {
		t.Errorf("Collection path mismatch. Got: %s, Want: %s", paths.Collection, expectedPaths.Collection)
	}
	if paths.Entities != expectedPaths.Entities {
		t.Errorf("Entities path mismatch. Got: %s, Want: %s", paths.Entities, expectedPaths.Entities)
	}
	if paths.Cache != expectedPaths.Cache {
		t.Errorf("Cache path mismatch. Got: %s, Want: %s", paths.Cache, expectedPaths.Cache)
	}
}
