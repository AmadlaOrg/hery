package validation

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntity(t *testing.T) {
	entityValidationService := NewEntityValidationService(&gitConfig.Config{})

	tests := []struct {
		name        string
		inputSchema *jsonschema.Schema
		heryContent map[string]any
		hasError    bool
	}{
		//
		// Error
		//
		{
			name:        "Error: _self should not be empty map",
			inputSchema: &jsonschema.Schema{},
			heryContent: map[string]any{
				"_type": "github.com/AmadlaOrg/Entity@latest",
				"_meta": map[string]any{
					"name":        "Entity",
					"description": "The root Entity definition.",
					"category":    "General",
					"tags": []any{
						"main",
						"master",
					},
				},
				"_self": map[string]any{},
			},
			hasError: true,
		},
		{
			name:        "Error: _self contains _type",
			inputSchema: &jsonschema.Schema{},
			heryContent: map[string]any{
				"_type": "github.com/AmadlaOrg/Entity@latest",
				"_meta": map[string]any{
					"name":        "Entity",
					"description": "The root Entity definition.",
					"category":    "General",
					"tags": []any{
						"main",
						"master",
					},
				},
				"_self": map[string]any{
					"_type":       "github.com/AmadlaOrg/Entity@latest",
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
