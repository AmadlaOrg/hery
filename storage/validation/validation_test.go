package validation

import (
	"github.com/AmadlaOrg/hery/storage"
	"testing"
)

// Mock function to simulate file existence
func mockFileExists(path string) bool {
	// Define which paths exist in your mock environment
	existingPaths := map[string]bool{
		"/home/user/.hery":                            true,
		"/home/user/.hery/collection":                 true,
		"/home/user/.hery/collection/test":            true,
		"/home/user/.hery/collection/test/entity":     true,
		"/home/user/.hery/collection/test/test.cache": true,
	}

	return existingPaths[path]
}

func Test_AllExist(t *testing.T) {
	// Backup the original fileExists function and restore it after the test
	originalFileExists := fileExists
	defer func() { fileExists = originalFileExists }()

	// Use the mock function for the test
	fileExists = mockFileExists

	// Define test cases
	tests := []struct {
		name      string
		paths     storage.AbsPaths
		expectErr bool
	}{
		{
			name: "All paths exist",
			paths: storage.AbsPaths{
				Storage:    "/home/user/.hery",
				Catalog:    "/home/user/.hery/collection",
				Collection: "/home/user/.hery/collection/test",
				Entities:   "/home/user/.hery/collection/test/entity",
				Cache:      "/home/user/.hery/collection/test/test.cache",
			},
			expectErr: false,
		},
		{
			name: "Storage path does not exist",
			paths: storage.AbsPaths{
				Storage:    "/home/user/.hery_not_exists",
				Catalog:    "/home/user/.hery/collection",
				Collection: "/home/user/.hery/collection/test",
				Entities:   "/home/user/.hery/collection/test/entity",
				Cache:      "/home/user/.hery/collection/test/test.cache",
			},
			expectErr: true,
		},
		{
			name: "Catalog path does not exist",
			paths: storage.AbsPaths{
				Storage:    "/home/user/.hery",
				Catalog:    "/home/user/.hery/collection_not_exists",
				Collection: "/home/user/.hery/collection/test",
				Entities:   "/home/user/.hery/collection/test/entity",
				Cache:      "/home/user/.hery/collection/test/test.cache",
			},
			expectErr: true,
		},
		{
			name: "Collection path does not exist",
			paths: storage.AbsPaths{
				Storage:    "/home/user/.hery",
				Catalog:    "/home/user/.hery/collection",
				Collection: "/home/user/.hery/collection/test_not_exists",
				Entities:   "/home/user/.hery/collection/test/entity",
				Cache:      "/home/user/.hery/collection/test/test.cache",
			},
			expectErr: true,
		},
		{
			name: "Entities path does not exist",
			paths: storage.AbsPaths{
				Storage:    "/home/user/.hery",
				Catalog:    "/home/user/.hery/collection",
				Collection: "/home/user/.hery/collection/test",
				Entities:   "/home/user/.hery/collection/test/entity_not_exists",
				Cache:      "/home/user/.hery/collection/test/test.cache",
			},
			expectErr: true,
		},
		{
			name: "Cache path does not exist",
			paths: storage.AbsPaths{
				Storage:    "/home/user/.hery",
				Catalog:    "/home/user/.hery/collection",
				Collection: "/home/user/.hery/collection/test",
				Entities:   "/home/user/.hery/collection/test/entity",
				Cache:      "/home/user/.hery/collection/test/test.cache_not_exists",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AllExist(tt.paths)
			if (err != nil) != tt.expectErr {
				t.Errorf("AllExist() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}
