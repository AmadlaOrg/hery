package validation

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntity(t *testing.T) {
	entityValidationService := New(&gitConfig.Config{})

	tests := []struct {
		name        string
		inputSchema *jsonschema.Schema
		heryContent map[string]any
		hasError    bool
	}{
		{
			name:        "Valid: minimal entity with _type",
			inputSchema: &jsonschema.Schema{},
			heryContent: map[string]any{
				"_type": "github.com/AmadlaOrg/Entity@latest",
			},
			hasError: false,
		},
		{
			name:        "Error: missing _type",
			inputSchema: &jsonschema.Schema{},
			heryContent: map[string]any{
				"_body": map[string]any{"key": "value"},
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
