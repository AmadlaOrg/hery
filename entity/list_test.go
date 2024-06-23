package entity

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageRoot(t *testing.T) {
	expectedDir := ""
	switch runtime.GOOS {
	case "windows":
		appDataDir := os.Getenv("APPDATA")
		expectedDir = filepath.Join(appDataDir, "Amadla", "entity")
	default: // "linux" and "darwin" (macOS)
		homeDir, err := os.UserHomeDir()
		assert.NoError(t, err)
		expectedDir = filepath.Join(homeDir, ".amadla", "entity")
	}

	dir, err := StorageRoot()
	assert.NoError(t, err)
	assert.Equal(t, expectedDir, dir)
}

func TestCrawlDirectoriesParallel(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "github.com", "AmadlaOrg")
	err := os.MkdirAll(filepath.Join(testDir, "Entity", "CPU@v1.0.0"), 0755)
	if err != nil {
		return
	}
	err = os.MkdirAll(filepath.Join(testDir, "Entity", "Memory@v1.0.0"), 0755)
	if err != nil {
		return
	}
	err = os.MkdirAll(filepath.Join(testDir, "Entity", "Net@v1.0.0"), 0755)
	if err != nil {
		return
	}
	err = os.MkdirAll(filepath.Join(testDir, "Entity", "Storage@v1.0.0"), 0755)
	if err != nil {
		return
	}

	entities, err := CrawlDirectoriesParallel(tmpDir)
	assert.NoError(t, err)
	expectedEntities := map[string]string{
		"CPU":     "1.0.0",
		"Memory":  "1.0.0",
		"Net":     "1.0.0",
		"Storage": "1.0.0",
	}

	assert.Equal(t, expectedEntities, entities)
}
