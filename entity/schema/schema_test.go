package schema

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindSchemaFile(t *testing.T) {
	s := &schemaImpl{}

	t.Run("finds single schema file", func(t *testing.T) {
		dir := t.TempDir()
		schemaPath := filepath.Join(dir, "application.hery.json")
		require.NoError(t, os.WriteFile(schemaPath, []byte(`{}`), 0644))

		result, err := s.FindSchemaFile(dir)
		assert.NoError(t, err)
		assert.Equal(t, schemaPath, result)
	})

	t.Run("error when no schema file found", func(t *testing.T) {
		dir := t.TempDir()

		_, err := s.FindSchemaFile(dir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no .hery.json schema file found")
	})

	t.Run("error when multiple schema files found", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "one.hery.json"), []byte(`{}`), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "two.hery.json"), []byte(`{}`), 0644))

		_, err := s.FindSchemaFile(dir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "multiple .hery.json schema files found")
	})

	t.Run("ignores .hery files (not .hery.json)", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "default.hery"), []byte("_type: test"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "package.hery.json"), []byte(`{}`), 0644))

		result, err := s.FindSchemaFile(dir)
		assert.NoError(t, err)
		assert.Equal(t, filepath.Join(dir, "package.hery.json"), result)
	})

	t.Run("glob error is propagated", func(t *testing.T) {
		origGlob := filepathGlob
		defer func() { filepathGlob = origGlob }()

		filepathGlob = func(pattern string) ([]string, error) {
			return nil, errors.New("glob failed")
		}

		_, err := s.FindSchemaFile("/some/path")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to search for schema file")
	})
}

func TestGenerateSchemaPath(t *testing.T) {
	s := &schemaImpl{}

	t.Run("returns found schema file when it exists", func(t *testing.T) {
		dir := t.TempDir()
		schemaPath := filepath.Join(dir, "template.hery.json")
		require.NoError(t, os.WriteFile(schemaPath, []byte(`{}`), 0644))

		result := s.GenerateSchemaPath(dir)
		assert.Equal(t, schemaPath, result)
	})

	t.Run("falls back to directory name when no schema file exists", func(t *testing.T) {
		dir := t.TempDir()
		// Rename to a known path so we can predict the fallback
		namedDir := filepath.Join(filepath.Dir(dir), "Application")
		os.Rename(dir, namedDir)
		defer os.Rename(namedDir, dir)

		result := s.GenerateSchemaPath(namedDir)
		assert.Equal(t, filepath.Join(namedDir, "application.hery.json"), result)
	})
}

func TestGenerateURN(t *testing.T) {
	s := &schemaImpl{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple entity URI",
			input:    "amadla.org/entity/application@v1.0.0",
			expected: "urn:hery:amadla.org:entity:application:v1.0.0",
		},
		{
			name:     "sub-type entity URI",
			input:    "amadla.org/entity/application/db@v1.0.0",
			expected: "urn:hery:amadla.org:entity:application:db:v1.0.0",
		},
		{
			name:     "github entity URI",
			input:    "github.com/AmadlaOrg/Application@v1.0.0",
			expected: "urn:hery:github.com:AmadlaOrg:Application:v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.GenerateURN(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
