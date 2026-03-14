package entity

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	versionPkg "github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/message"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/stretchr/testify/assert"
)

func TestSetContent(t *testing.T) {
	s := &service{}

	tests := []struct {
		name        string
		entity      Entity
		input       NotFormatedContent
		expectType  string
		expectError bool
	}{
		{
			name:   "Full content",
			entity: Entity{},
			input: NotFormatedContent{
				"_type":   "github.com/AmadlaOrg/Entity@v1.0.0",
				"_extends": "parent-entity",
				"_meta":   map[string]any{"name": "Test"},
				"_body":   map[string]any{"key": "value"},
			},
			expectType: "github.com/AmadlaOrg/Entity@v1.0.0",
		},
		{
			name:   "Type from entity URI",
			entity: Entity{Uri: "github.com/AmadlaOrg/Entity@v2.0.0"},
			input: NotFormatedContent{
				"_type": "github.com/AmadlaOrg/Entity@v1.0.0",
			},
			expectType: "github.com/AmadlaOrg/Entity@v2.0.0",
		},
		{
			name:   "Missing _type",
			entity: Entity{},
			input: NotFormatedContent{
				"_body": map[string]any{"key": "value"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := s.setContent(tt.entity, tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectType, content.Type)
			}
		})
	}
}

func TestFindDir(t *testing.T) {
	basePath := filepath.Join(t.TempDir(), "entity")

	// Create test directories
	exactDir := filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApp@v1.0.0")
	err := os.MkdirAll(exactDir, os.ModePerm)
	assert.NoError(t, err)

	pseudoDir := filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApp@v0.0.0-20240726095222-c7e9911d38b2")
	err = os.MkdirAll(pseudoDir, os.ModePerm)
	assert.NoError(t, err)

	mockVersionValidation := &versionValidationPkg.MockEntityVersionValidation{}
	mockVersion := &versionPkg.MockEntityVersion{}
	s := &service{
		EntityVersion:           mockVersion,
		EntityVersionValidation: mockVersionValidation,
	}

	tests := []struct {
		name        string
		paths       storage.AbsPaths
		entityVals  Entity
		isPseudo    bool
		expected    string
		expectError bool
	}{
		{
			name:  "Exact version match",
			paths: storage.AbsPaths{Entities: basePath},
			entityVals: Entity{
				Name:    "EntityApp",
				Version: "v1.0.0",
				Origin:  "github.com/AmadlaOrg",
				Uri:     exactDir,
			},
			isPseudo: false,
			expected: exactDir,
		},
		{
			name:  "Exact version not found",
			paths: storage.AbsPaths{Entities: basePath},
			entityVals: Entity{
				Name:    "EntityApp",
				Version: "v9.9.9",
				Origin:  "github.com/AmadlaOrg",
				Uri:     filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApp@v9.9.9"),
			},
			isPseudo:    false,
			expectError: true,
		},
		{
			name:  "Pseudo version match",
			paths: storage.AbsPaths{Entities: basePath},
			entityVals: Entity{
				Name:    "EntityApp",
				Version: "v0.0.0-20240726095222-c7e9911d38b2",
				Origin:  "github.com/AmadlaOrg",
			},
			isPseudo: true,
			expected: pseudoDir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVersionValidation.ExpectedCalls = nil
			mockVersion.ExpectedCalls = nil
			mockVersionValidation.EXPECT().PseudoFormat(tt.entityVals.Version).Return(tt.isPseudo)
			if tt.isPseudo {
				mockVersion.EXPECT().GeneratePseudoPattern(tt.entityVals.Name, tt.entityVals.Version).
					Return(fmt.Sprintf("%s@%s-*-%s", tt.entityVals.Name, tt.entityVals.Version[:6], tt.entityVals.Version[22:]))
			}

			result, err := s.FindDir(tt.paths, tt.entityVals)
			if tt.expectError {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, message.ErrorNotFound) || errors.Is(err, message.ErrorMultipleFound) || err != nil)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCheckDuplicate(t *testing.T) {
	mockVersionValidation := &versionValidationPkg.MockEntityVersionValidation{}
	mockVersion := &versionPkg.MockEntityVersion{}
	s := &service{
		EntityVersion:           mockVersion,
		EntityVersionValidation: mockVersionValidation,
	}

	tests := []struct {
		name       string
		entities   []Entity
		entityMeta Entity
		isPseudo   bool
		expectErr  bool
	}{
		{
			name: "Exact version duplicate",
			entities: []Entity{
				{Name: "EntityApp", Version: "v1.0.0", Origin: "github.com/AmadlaOrg"},
			},
			entityMeta: Entity{Name: "EntityApp", Version: "v1.0.0", Origin: "github.com/AmadlaOrg"},
			isPseudo:   false,
			expectErr:  true,
		},
		{
			name: "Different version no duplicate",
			entities: []Entity{
				{Name: "EntityApp", Version: "v1.0.0", Origin: "github.com/AmadlaOrg"},
			},
			entityMeta: Entity{Name: "EntityApp", Version: "v2.0.0", Origin: "github.com/AmadlaOrg"},
			isPseudo:   false,
			expectErr:  false,
		},
		{
			name:       "Empty entities list",
			entities:   []Entity{},
			entityMeta: Entity{Name: "EntityApp", Version: "v1.0.0", Origin: "github.com/AmadlaOrg"},
			isPseudo:   false,
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVersionValidation.ExpectedCalls = nil
			mockVersionValidation.EXPECT().PseudoFormat(tt.entityMeta.Version).Return(tt.isPseudo).Maybe()
			for _, e := range tt.entities {
				mockVersionValidation.EXPECT().PseudoFormat(e.Version).Return(tt.isPseudo).Maybe()
			}

			err := s.CheckDuplicate(tt.entities, tt.entityMeta)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "duplicate")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGeneratePseudoVersionPattern(t *testing.T) {
	entityService := New(&gitConfig.Config{})

	tests := []struct {
		name         string
		inputName    string
		inputVersion string
		expected     string
	}{
		{
			name:         "Basic pseudo version",
			inputName:    "EntityApplication",
			inputVersion: "v0.0.0-20231231235959-1234567890ab",
			expected:     "EntityApplication@v0.0.0-*-1234567890ab",
		},
		{
			name:         "Another pseudo version",
			inputName:    "AnotherEntity",
			inputVersion: "v1.2.3-20230101010101-abcdefabcdef",
			expected:     "AnotherEntity@v1.2.3-*-abcdefabcdef",
		},
		{
			name:         "Different pseudo version format",
			inputName:    "ThirdEntity",
			inputVersion: "v2.0.0-20231231235959-1234567890ab",
			expected:     "ThirdEntity@v2.0.0-*-1234567890ab",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := entityService.GeneratePseudoVersionPattern(test.inputName, test.inputVersion)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestReadAll_Unit(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test .hery file
	content := []byte("_type: example.com/Entity@v1.0.0\n_body:\n  name: test\n")
	err := os.WriteFile(filepath.Join(tmpDir, "test.hery"), content, 0644)
	assert.NoError(t, err)

	// Create a non-.hery file (should be ignored)
	err = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("ignore me"), 0644)
	assert.NoError(t, err)

	s := &service{}
	docs, err := s.ReadAll(tmpDir)
	assert.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, "example.com/Entity@v1.0.0", docs[0]["_type"])
}

func TestReadAll_Unit_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	s := &service{}
	docs, err := s.ReadAll(tmpDir)
	assert.NoError(t, err)
	assert.Empty(t, docs)
}

func TestReadAll_Unit_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	content := []byte("invalid: yaml: [broken")
	err := os.WriteFile(filepath.Join(tmpDir, "bad.hery"), content, 0644)
	assert.NoError(t, err)

	s := &service{}
	_, err = s.ReadAll(tmpDir)
	assert.Error(t, err)
}
