package heryext

import (
	"testing"
)

func TestRead(t *testing.T) {
	//heryExtService := NewHeryExtService()
	//tmpDir := t.TempDir()

	/*t.Run("read .yml file", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "config.yml")
		content := []byte(`key: value`)
		err := os.WriteFile(filePath, content, 0644)
		assert.NoError(t, err)

		data, err := heryExtService.Read(tmpDir, "config")
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

		data, err := heryExtService.Read(tmpDir, "config")
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

		_, err = heryExtService.Read(tmpDir, "config")
		assert.Error(t, err)
		assert.Equal(t, "both "+ymlFilePath+", "+yamlFilePath+" exists", err.Error())
	})

	t.Run("file does not exist", func(t *testing.T) {
		_, err := heryExtService.Read(tmpDir, "nonexistent")
		assert.Error(t, err)
		assert.Equal(t, filepath.Join(tmpDir, "nonexistent.yml")+" does not exist", err.Error())
	})*/

	/*t.Run("error reading the file", func(t *testing.T) {
		// Create a file with no read permissions
		filePath := filepath.Join(tmpDir, "config_no_read.yml")
		content := []byte(`key: value`)
		err := os.WriteFile(filePath, content, 0000)
		assert.NoError(t, err)
		defer func(name string, mode os.FileMode) {
			err := os.Chmod(name, mode)
			if err != nil {
				assert.FailNow(t, err.Error())
			}
		}(filePath, 0644) // Ensure the file can be removed after the test

		_, err = heryExtService.Read(tmpDir, "config_no_read")
		assert.Error(t, err)
	})

	t.Run("error unmarshaling the content", func(t *testing.T) {
		// Create a file with invalid YAML content
		filePath := filepath.Join(tmpDir, "invalid_config.yml")
		content := []byte(`: key value :`)
		err := os.WriteFile(filePath, content, 0644)
		assert.NoError(t, err)

		_, err = heryExtService.Read(tmpDir, "invalid_config")
		assert.Error(t, err)
	})*/
}
