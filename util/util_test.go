package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("file exists", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "file.txt")
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		assert.NoError(t, err)

		exists := FileExists(filePath)
		assert.True(t, exists)
	})

	t.Run("file does not exist", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "nonexistent.txt")

		exists := FileExists(filePath)
		assert.False(t, exists)
	})

	t.Run("directory exists", func(t *testing.T) {
		dirPath := filepath.Join(tmpDir, "subdir")
		err := os.Mkdir(dirPath, 0755)
		assert.NoError(t, err)

		exists := FileExists(dirPath)
		assert.True(t, exists)
	})
}

func TestReadYaml(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("read .yml file", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "config.yml")
		content := []byte(`key: value`)
		err := os.WriteFile(filePath, content, 0644)
		assert.NoError(t, err)

		data, err := ReadYaml(tmpDir, "config")
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"key": "value"}, data)
	})

	t.Run("read .yaml file", func(t *testing.T) {
		// Ensure no .yml file exists
		ymlFilePath := filepath.Join(tmpDir, "config.yml")
		_ = os.Remove(ymlFilePath)

		filePath := filepath.Join(tmpDir, "config.yaml")
		content := []byte(`key: value`)
		err := os.WriteFile(filePath, content, 0644)
		assert.NoError(t, err)

		data, err := ReadYaml(tmpDir, "config")
		assert.NoError(t, err)
		assert.Equal(t, map[string]interface{}{"key": "value"}, data)
	})

	t.Run("both .yml and .yaml files exist", func(t *testing.T) {
		ymlFilePath := filepath.Join(tmpDir, "config.yml")
		yamlFilePath := filepath.Join(tmpDir, "config.yaml")
		content := []byte(`key: value`)
		err := os.WriteFile(ymlFilePath, content, 0644)
		assert.NoError(t, err)
		err = os.WriteFile(yamlFilePath, content, 0644)
		assert.NoError(t, err)

		_, err = ReadYaml(tmpDir, "config")
		assert.Error(t, err)
		assert.Equal(t, "both "+ymlFilePath+", "+yamlFilePath+" exists", err.Error())
	})

	t.Run("file does not exist", func(t *testing.T) {
		_, err := ReadYaml(tmpDir, "nonexistent")
		assert.Error(t, err)
		assert.Equal(t, filepath.Join(tmpDir, "nonexistent.yml")+" does not exist", err.Error())
	})
}
