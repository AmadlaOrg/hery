package validation

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestEntity(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "test", "fixture")
	validEntityAbsPath, err := filepath.Abs(filepath.Join(fixturePath, "/valid-entity"))
	if err != nil {
		t.Fatal(err)
	}

	entityValidationService := NewEntityValidationService()

	tests := []struct {
		name                string
		inputEntityPath     string
		inputCollectionName string
		heryContent         map[string]any
		hasError            bool
	}{
		{
			name:                "valid",
			inputEntityPath:     validEntityAbsPath,
			inputCollectionName: "amadla",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = entityValidationService.Entity(tt.inputEntityPath, tt.inputCollectionName, tt.heryContent)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
