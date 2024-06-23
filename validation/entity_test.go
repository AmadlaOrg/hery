package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	tmpDir := t.TempDir()

	// Create schema file
	schemaDir := filepath.Join(tmpDir, ".amadla")
	err := os.MkdirAll(schemaDir, 0755)
	assert.NoError(t, err)
	schemaContent := `{
		"type": "object",
		"properties": {
			"key": { "type": "string" }
		},
		"required": ["key"]
	}`
	err = os.WriteFile(filepath.Join(schemaDir, "schema.json"), []byte(schemaContent), 0644)
	assert.NoError(t, err)

	t.Run("valid YAML", func(t *testing.T) {
		// Create a valid YAML file
		yamlContent := `key: value`
		err := os.WriteFile(filepath.Join(tmpDir, "entity.yaml"), []byte(yamlContent), 0644)
		assert.NoError(t, err)

		// Validate the entity
		err = Entity(tmpDir)
		assert.NoError(t, err)
	})

	t.Run("invalid YAML", func(t *testing.T) {
		// Create an invalid YAML file
		yamlContent := `invalidKey: value`
		err := os.WriteFile(filepath.Join(tmpDir, "entity.yaml"), []byte(yamlContent), 0644)
		assert.NoError(t, err)

		// Validate the entity
		err = Entity(tmpDir)
		assert.Error(t, err)
	})

	t.Run("YAML file not found", func(t *testing.T) {
		// Remove any existing YAML files
		_ = os.Remove(filepath.Join(tmpDir, "entity.yaml"))
		_ = os.Remove(filepath.Join(tmpDir, "entity.yml"))

		// Validate the entity
		err := Entity(tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "YAML file not found in entity")
	})

	t.Run("schema file not found", func(t *testing.T) {
		// Create a valid YAML file
		yamlContent := `key: value`
		err := os.WriteFile(filepath.Join(tmpDir, "entity.yaml"), []byte(yamlContent), 0644)
		assert.NoError(t, err)

		// Remove the schema file
		_ = os.Remove(filepath.Join(schemaDir, "schema.json"))

		// Validate the entity
		err = Entity(tmpDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error loading JSON schema")
	})
}
