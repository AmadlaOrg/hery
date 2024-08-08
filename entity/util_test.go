package entity

import (
	"fmt"
	"github.com/AmadlaOrg/hery/storage"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindEntityDir(t *testing.T) {
	// Setup test directories
	basePath := "/tmp/.hery/test/entity"
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		t.Fatal("cannot create test directory")
	}
	defer func() {
		err := os.RemoveAll("/tmp/.hery")
		if err != nil {
			t.Fatal("cannot remove test directory")
		}
	}() // Clean up after tests

	tests := []struct {
		name        string
		paths       storage.AbsPaths
		entityVals  Entity
		setupFunc   func()
		expected    string
		expectedErr error
	}{
		{
			name: "Exact version match",
			paths: storage.AbsPaths{
				Entities: basePath,
			},
			entityVals: Entity{
				Name:    "EntityApplication",
				Version: "v0.0.0",
				Origin:  "github.com/AmadlaOrg",
				Entity:  filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0"),
			},
			setupFunc: func() {
				err := os.MkdirAll(filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0"), os.ModePerm)
				if err != nil {
					t.Fatal("cannot create test directory")
				}
			},
			expected:    filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0"),
			expectedErr: nil,
		},
		{
			name: "Pseudo version match",
			paths: storage.AbsPaths{
				Entities: basePath,
			},
			entityVals: Entity{
				Name:    "EntityApplication",
				Version: "v0.0.0-20240726095222-c7e9911d38b2",
				Origin:  "github.com/AmadlaOrg",
			},
			setupFunc: func() {
				err := os.MkdirAll(filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0-20240726095222-c7e9911d38b2"), os.ModePerm)
				if err != nil {
					t.Fatal("cannot create test directory")
				}
			},
			expected:    filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0-20240726095222-c7e9911d38b2"),
			expectedErr: nil,
		},
		/*{
			name: "No matching exact version",
			paths: storage.AbsPaths{
				Entities: basePath,
			},
			entityVals: Entity{
				Name:    "EntityApplication",
				Version: "v0.0.0",
				Origin:  "github.com/AmadlaOrg",
				Entity:  filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0"),
			},
			setupFunc: func() {},
			expected:  "",
			expectedErr: errors.Join(
				errtypes.NotFoundError,
				fmt.Errorf("no matching directory found for exact version: %s", filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0")),
			),
		},
		{
			name: "No matching pseudo version",
			paths: storage.AbsPaths{
				Entities: basePath,
			},
			entityVals: Entity{
				Name:    "EntityApplication",
				Version: "v0.0.0-20240726095222-c7e9911d38b2",
				Origin:  "github.com/AmadlaOrg",
			},
			setupFunc: func() {},
			expected:  "",
			expectedErr: errors.Join(
				errtypes.NotFoundError,
				fmt.Errorf("no matching directories found for pattern: %s", filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0-*c7e9911d38b2")),
			),
		},
		{
			name: "Multiple matching pseudo versions",
			paths: storage.AbsPaths{
				Entities: basePath,
			},
			entityVals: Entity{
				Name:    "EntityApplication",
				Version: "v0.0.0-20240726095222-c7e9911d38b2",
				Origin:  "github.com/AmadlaOrg",
			},
			setupFunc: func() {
				err := os.MkdirAll(filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0-20240726095222-c7e9911d38b2"), os.ModePerm)
				if err != nil {
					t.Fatal("cannot create test directory")
				}
				err = os.MkdirAll(filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0-20240726095322-c7e9911d38b2"), os.ModePerm)
				if err != nil {
					t.Fatal("cannot create test directory")
				}
			},
			expected: "",
			expectedErr: errors.Join(
				errtypes.MultipleFoundError,
				fmt.Errorf("multiple matching directories found for pattern: %s", filepath.Join(basePath, "github.com/AmadlaOrg", "EntityApplication@v0.0.0-*c7e9911d38b2")),
			),
		},*/
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup test case
			test.setupFunc()

			result, err := FindEntityDir(test.paths, test.entityVals)
			if test.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expectedErr.Error())
			}
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestCheckDuplicateEntity(t *testing.T) {
	tests := []struct {
		name       string
		entities   []Entity
		entityMeta Entity
		expected   error
	}{
		{
			name: "Exact version match",
			entities: []Entity{
				{
					Name:    "EntityApplication",
					Version: "v0.0.0",
					Origin:  "github.com/AmadlaOrg",
				},
			},
			entityMeta: Entity{
				Name:    "EntityApplication",
				Version: "v0.0.0",
				Origin:  "github.com/AmadlaOrg",
			},
			expected: fmt.Errorf("duplicate entity found: %v", Entity{
				Name:    "EntityApplication",
				Version: "v0.0.0",
				Origin:  "github.com/AmadlaOrg",
			}),
		},
		{
			name: "Pseudo version match",
			entities: []Entity{
				{
					Name:    "EntityApplication",
					Version: "v0.0.0-20240726095222-c7e9911d38b2",
					Origin:  "github.com/AmadlaOrg",
				},
			},
			entityMeta: Entity{
				Name:    "EntityApplication",
				Version: "v0.0.0-20240726095222-c7e9911d38b2",
				Origin:  "github.com/AmadlaOrg",
			},
			expected: fmt.Errorf("duplicate entity found: %v", Entity{
				Name:    "EntityApplication",
				Version: "v0.0.0-20240726095222-c7e9911d38b2",
				Origin:  "github.com/AmadlaOrg",
			}),
		},
		{
			name: "No match",
			entities: []Entity{
				{
					Name:    "EntityApplication",
					Version: "v0.0.0",
					Origin:  "github.com/AmadlaOrg",
				},
			},
			entityMeta: Entity{
				Name:    "EntityApplication",
				Version: "v0.0.1",
				Origin:  "github.com/AmadlaOrg",
			},
			expected: nil,
		},
		{
			name: "Pseudo version no match",
			entities: []Entity{
				{
					Name:    "EntityApplication",
					Version: "v0.0.0-20240726095222-c7e9911d38b2",
					Origin:  "github.com/AmadlaOrg",
				},
			},
			entityMeta: Entity{
				Name:    "EntityApplication",
				Version: "v0.0.0-20240726095222-c889911d00b2",
				Origin:  "github.com/AmadlaOrg",
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckDuplicateEntity(test.entities, test.entityMeta)
			if test.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expected.Error())
			}
		})
	}
}

func TestGeneratePseudoVersionPattern(t *testing.T) {
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
			result := GeneratePseudoVersionPattern(test.inputName, test.inputVersion)
			assert.Equal(t, test.expected, result)
		})
	}
}