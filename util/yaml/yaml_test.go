package yaml

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("read .yml file", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "config.yml")
		content := []byte(`key: value`)
		err := os.WriteFile(filePath, content, 0644)
		assert.NoError(t, err)

		data, err := Read(tmpDir, "config")
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

		data, err := Read(tmpDir, "config")
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

		_, err = Read(tmpDir, "config")
		assert.Error(t, err)
		assert.Equal(t, "both "+ymlFilePath+", "+yamlFilePath+" exists", err.Error())
	})

	t.Run("file does not exist", func(t *testing.T) {
		_, err := Read(tmpDir, "nonexistent")
		assert.Error(t, err)
		assert.Equal(t, filepath.Join(tmpDir, "nonexistent.yml")+" does not exist", err.Error())
	})
}
