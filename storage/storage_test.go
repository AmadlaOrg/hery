package storage

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestStorageRoot(t *testing.T) {
	expectedDir := ""
	switch runtime.GOOS {
	case "windows":
		appDataDir := os.Getenv("APPDATA")
		expectedDir = filepath.Join(appDataDir, "Hery", "entity")
	default: // "linux" and "darwin" (macOS)
		homeDir, err := os.UserHomeDir()
		assert.NoError(t, err)
		expectedDir = filepath.Join(homeDir, ".hery", "entity")
	}

	dir, err := Path()
	assert.NoError(t, err)
	assert.Equal(t, expectedDir, dir)
}
