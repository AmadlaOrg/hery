package validation

import (
	"github.com/AmadlaOrg/hery/storage"
	"testing"
)

func mockFileExists(path string) bool {
	existingPaths := map[string]bool{
		"/home/user/.cache/hery":        true,
		"/home/user/.cache/hery/entity": true,
	}
	return existingPaths[path]
}

func Test_AllExist(t *testing.T) {
	originalFileExists := fileExists
	defer func() { fileExists = originalFileExists }()

	fileExists = mockFileExists

	tests := []struct {
		name      string
		paths     storage.AbsPaths
		expectErr bool
	}{
		{
			name: "All paths exist",
			paths: storage.AbsPaths{
				Storage:  "/home/user/.cache/hery",
				Entities: "/home/user/.cache/hery/entity",
			},
			expectErr: false,
		},
		{
			name: "Storage path does not exist",
			paths: storage.AbsPaths{
				Storage:  "/home/user/.cache/hery_not_exists",
				Entities: "/home/user/.cache/hery/entity",
			},
			expectErr: true,
		},
		{
			name: "Entities path does not exist",
			paths: storage.AbsPaths{
				Storage:  "/home/user/.cache/hery",
				Entities: "/home/user/.cache/hery/entity_not_exists",
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
