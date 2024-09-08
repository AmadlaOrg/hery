package schema

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestLoadSchemaFile(t *testing.T) {
	baseSchemaAbsPath, err := filepath.Abs(filepath.Join("..", "..", ".schema", "entity.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	schemaService := SSchema{}
	tests := []struct {
		name            string
		inputSchemaPath string
		expected        map[string]any
		hasError        bool
	}{
		{
			name:            "valid",
			inputSchemaPath: baseSchemaAbsPath,
			expected: map[string]any{
				"title": "HERY basic configuration properties.",
				"type":  "object",
				"properties": map[string]any{
					"_entity": map[string]any{
						"type":        "string",
						"format":      "uri",
						"description": "The URI that uniquely identifies the entity within the HERY system.",
					},
					"_id": map[string]any{
						"type":        "string",
						"pattern":     "^[a-zA-Z0-9_\\-:/]+$",
						"description": "A unique identifier for an entity dataset.",
					},
					"_self": map[string]any{
						"type":                 "object",
						"description":          "Used to reference the current entity inside another entity, so that there is no need to use _entity to define the entity block. In other words, it is a shorthand for _entity to reference the current entity.",
						"additionalProperties": true,
					},
				},
				"additionalProperties": false,
				"$schema":              "https://json-schema.org/draft/2020-12/schema",
				"id":                   "https://raw.githubusercontent.com/AmadlaOrg/hery/master/.schema/entity.schema.json",
				"description":          "Defines the foundational HERY entity properties. This schema sets the foundational structure for entities in the HERY system, including reserved properties that are crucial for entity identification and referencing within the system.",
				"required": []any{
					"_entity",
				},
				"propertiesPattern": map[string]any{
					"^[^_].*$": map[string]any{
						"type": "any",
					},
				},
			},
			hasError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := schemaService.loadSchemaFile(test.inputSchemaPath)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected, got)
		})
	}
}

func TestMergeSchemas(t *testing.T) {
	schemaService := SSchema{}
	tests := []struct {
		name            string
		inputBaseSchema map[string]any
		inputMainSchema map[string]any
		expected        map[string]any
	}{
		{
			name: "valid",
			inputBaseSchema: map[string]any{
				"title": "HERY basic configuration properties.",
				"type":  "object",
				"properties": map[string]any{
					"_entity": map[string]any{
						"type":        "string",
						"format":      "uri",
						"description": "The URI that uniquely identifies the entity within the HERY system.",
					},
					"_id": map[string]any{
						"type":        "string",
						"pattern":     "^[a-zA-Z0-9_\\-:/]+$",
						"description": "A unique identifier for an entity dataset.",
					},
					"_self": map[string]any{
						"type":                 "object",
						"description":          "Used to reference the current entity inside another entity, so that there is no need to use _entity to define the entity block. In other words, it is a shorthand for _entity to reference the current entity.",
						"additionalProperties": true,
					},
				},
				"additionalProperties": false,
				"$schema":              "https://json-schema.org/draft/2020-12/schema",
				"id":                   "https://raw.githubusercontent.com/AmadlaOrg/hery/master/.schema/entity.schema.json",
				"description":          "Defines the foundational HERY entity properties. This schema sets the foundational structure for entities in the HERY system, including reserved properties that are crucial for entity identification and referencing within the system.",
				"required": []any{
					"_entity",
				},
				"propertiesPattern": map[string]any{
					"^[^_].*$": map[string]any{
						"type": "any",
					},
				},
			},
			inputMainSchema: map[string]any{},
			expected:        map[string]any{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := schemaService.mergeSchemas(test.inputBaseSchema, test.inputMainSchema)
			assert.Equal(t, test.expected, got)
		})
	}
}

/*func TestLoadJSONSchema(t *testing.T) {
	schemaService := NewSchemaService()
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
		schema, err := schemaService.Load(schemaPath)
		assert.NoError(t, err)
		assert.NotNil(t, schema)
	})

	t.Run("schema file does not exist", func(t *testing.T) {
		// Attempt to load a non-existent schema file
		schemaPath := filepath.Join(tmpDir, "nonexistent.json")

		// Load schema and check
		schema, err := schemaService.Load(schemaPath)
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
		schema, err := schemaService.Load(schemaPath)
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
		schema, err := schemaService.Load(schemaPath)
		assert.Error(t, err)
		assert.Nil(t, schema)
	})
}*/
