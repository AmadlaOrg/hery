package heryext

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	fixturePath := filepath.Join("..", "test", "fixture")
	validEntityAbsPath, err := filepath.Abs(filepath.Join(fixturePath, "/valid-entity"))
	if err != nil {
		t.Fatal(err)
	}

	heryExtService := NewHeryExtService()
	tests := []struct {
		name                string
		inputPath           string
		inputCollectionName string
		expected            map[string]any
		hasError            bool
	}{
		{
			name:                "Valid",
			inputPath:           validEntityAbsPath,
			inputCollectionName: "amadla",
			expected: map[string]any{
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
		//
		// Error
		//
		{
			name:                "Error: file not found",
			inputPath:           fixturePath,
			inputCollectionName: "amadla",
			expected:            nil,
			hasError:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			heryContent, err := heryExtService.Read(tt.inputPath, tt.inputCollectionName)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, heryContent)
		})
	}
}
