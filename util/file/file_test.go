package file

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("file exists", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "file.txt")
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		assert.NoError(t, err)

		exists := Exists(filePath)
		assert.True(t, exists)
	})

	t.Run("file does not exist", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "nonexistent.txt")

		exists := Exists(filePath)
		assert.False(t, exists)
	})

	t.Run("directory exists", func(t *testing.T) {
		dirPath := filepath.Join(tmpDir, "subdir")
		err := os.Mkdir(dirPath, 0755)
		assert.NoError(t, err)

		exists := Exists(dirPath)
		assert.True(t, exists)
	})
}
