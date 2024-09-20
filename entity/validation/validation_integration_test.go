package validation

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntity(t *testing.T) {
	/*fixturePath := filepath.Join("..", "..", "test", "fixture")
	validEntityAbsPath, err := filepath.Abs(filepath.Join(fixturePath, "/valid-entity"))
	if err != nil {
		t.Fatal(err)
	}*/

	entityValidationService := NewEntityValidationService()

	tests := []struct {
		name                string
		inputCollectionName string
		inputSchema         *jsonschema.Schema
		heryContent         map[string]any
		hasError            bool
	}{
		{
			name:                "Valid: Same entityUri entity with root of `.hery` without `_self`",
			inputCollectionName: "amadla",
			inputSchema:         &jsonschema.Schema{},
			heryContent: map[string]any{
				"_entity":     "github.com/AmadlaOrg/Entity@latest",
				"name":        "Entity",
				"description": "The root Entity definition.",
				"category":    "General",
				"tags": []any{
					"main",
					"master",
				},
			},
			hasError: false,
		},
		{
			name:                "Valid: With `_self`",
			inputCollectionName: "amadla",
			inputSchema:         &jsonschema.Schema{},
			//inputEntityUri:      "github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion@latest",
			heryContent: map[string]any{
				"_entity":     "github.com/AmadlaOrg/Entity@latest",
				"name":        "Entity",
				"description": "The root Entity definition.",
				"category":    "General",
				"tags": []any{
					"main",
					"master",
				},
				"_self": map[string]any{
					"title":       "Some random title",
					"description": "The random title description.",
				},
			},
			hasError: false,
		},
		//
		// Error
		//
		{
			name:                "Error: entityUri should not be of the same entity that is set in the root of `.hery` if _self is set",
			inputCollectionName: "amadla",
			inputSchema:         &jsonschema.Schema{},
			//inputEntityUri:      "github.com/AmadlaOrg/Entity@latest",
			heryContent: map[string]any{
				"_entity":     "github.com/AmadlaOrg/Entity@latest",
				"name":        "Entity",
				"description": "The root Entity definition.",
				"category":    "General",
				"tags": []any{
					"main",
					"master",
				},
				"_self": map[string]any{},
			},
			hasError: true,
		},
		{
			name:                "Error: With `_self` that contains `_entity`",
			inputSchema:         &jsonschema.Schema{},
			inputCollectionName: "amadla",
			//inputEntityUri:      "github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion@latest",
			heryContent: map[string]any{
				"_entity":     "github.com/AmadlaOrg/Entity@latest",
				"name":        "Entity",
				"description": "The root Entity definition.",
				"category":    "General",
				"tags": []any{
					"main",
					"master",
				},
				"_self": map[string]any{
					"_entity":     "github.com/AmadlaOrg/Entity@latest",
					"title":       "Some random title",
					"description": "The random title description.",
				},
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := entityValidationService.Entity(
				tt.inputCollectionName,
				tt.inputSchema,
				tt.heryContent)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
