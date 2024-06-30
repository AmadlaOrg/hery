package validation

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestEntity(t *testing.T) {
	tmpDir := t.TempDir()

	// Create schema file
	schemaDir := filepath.Join(tmpDir, ".hery")
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

func TestEntityUrl(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"github.com/user/module", true},
		{"github.com/user/module@v1.0.0", true},
		{"https://github.com/user/module", false},
		{"github.com/user/module with spaces", false},
		{"github.com/user/module:colon", false},
		{"github.com/user/module?query", false},
		{"github.com/user/module&param", false},
		{"github.com/user/module=equal", false},
		{"github.com/user/module#fragment", false},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			result := EntityUrl(test.path)
			if result != test.expected {
				t.Errorf("EntityUrl(%q) = %v; expected %v", test.path, result, test.expected)
			}
		})
	}
}
