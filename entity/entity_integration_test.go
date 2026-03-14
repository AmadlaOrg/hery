package entity

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestReadAll(t *testing.T) {
	fixturePath := filepath.Join("..", "test", "fixture")
	validEntityAbsPath, err := filepath.Abs(filepath.Join(fixturePath, "/valid-entity"))
	if err != nil {
		t.Fatal(err)
	}

	heryExtService := New(&gitConfig.Config{})
	tests := []struct {
		name      string
		inputPath string
		hasError  bool
	}{
		{
			name:      "Valid",
			inputPath: validEntityAbsPath,
			hasError:  false,
		},
		//
		// Error
		//
		{
			name:      "Error: directory not found",
			inputPath: "/nonexistent/path",
			hasError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := heryExtService.ReadAll(tt.inputPath)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, results)
			}
		})
	}
}
