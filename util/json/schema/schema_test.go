package schema

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadJSONSchema(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("load valid schema", func(t *testing.T) {
		// Create a valid schema file
		schemaContent := `{
			"type": "object",
			"properties": {
				"key": { "type": "string" }
			},
			"required": ["key"]
		}`
		schemaPath := filepath.Join(tmpDir, "schema.json")
		err := os.WriteFile(schemaPath, []byte(schemaContent), 0644)
		assert.NoError(t, err)

		// Load schema and check
		schema, err := Load(schemaPath)
		assert.NoError(t, err)
		assert.NotNil(t, schema)
	})

	t.Run("schema file does not exist", func(t *testing.T) {
		// Attempt to load a non-existent schema file
		schemaPath := filepath.Join(tmpDir, "nonexistent.json")

		// Load schema and check
		schema, err := Load(schemaPath)
		assert.Error(t, err)
		assert.Nil(t, schema)
	})

	t.Run("load invalid schema", func(t *testing.T) {
		// Create an invalid schema file
		schemaContent := `{
			"type": "object",
			"properties": {
				"key": { "type": "unknown" }
			}
		}`
		schemaPath := filepath.Join(tmpDir, "invalid_schema.json")
		err := os.WriteFile(schemaPath, []byte(schemaContent), 0644)
		assert.NoError(t, err)

		// Load schema and check
		schema, err := Load(schemaPath)
		assert.Error(t, err)
		assert.Nil(t, schema)
	})

	t.Run("add resource failure", func(t *testing.T) {
		// Create a schema file with no read permissions
		schemaPath := filepath.Join(tmpDir, "unreadable_schema.json")
		err := os.WriteFile(schemaPath, []byte(`{}`), 0000) // No read permissions
		assert.NoError(t, err)
		defer func(name string, mode os.FileMode) {
			err := os.Chmod(name, mode)
			if err != nil {
				assert.FailNow(t, err.Error())
			}
		}(schemaPath, 0644) // Ensure the file can be removed after the test

		// Load schema and check
		schema, err := Load(schemaPath)
		assert.Error(t, err)
		assert.Nil(t, schema)
	})
}
