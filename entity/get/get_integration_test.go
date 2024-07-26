package get

import (
	"github.com/AmadlaOrg/hery/storage"
	"os"
	"path/filepath"
	"testing"
)

func Test_Integration_Get(t *testing.T) {
	tests := []struct {
		name      string
		paths     storage.AbsPaths
		entityUri []string
	}{
		{
			name: "Get One",
			paths: storage.AbsPaths{
				Storage:    filepath.Join(os.TempDir(), ".hery"),
				Collection: filepath.Join(os.TempDir(), ".hery", "collection"),
				Entities:   filepath.Join(os.TempDir(), ".hery", "collection", "testone", "entity"),
				Cache:      filepath.Join(os.TempDir(), ".hery", "collection", "testone", "test.cache"),
			},
			entityUri: []string{
				"github.com/AmadlaOrg/EntityApplication",
			},
		},
		{
			name: "Get Multiple Paths",
			paths: storage.AbsPaths{
				Storage:    filepath.Join(os.TempDir(), ".hery"),
				Collection: filepath.Join(os.TempDir(), ".hery", "collection"),
				Entities:   filepath.Join(os.TempDir(), ".hery", "collection", "testmultiple", "entity"),
				Cache:      filepath.Join(os.TempDir(), ".hery", "collection", "testmultiple", "test.cache"),
			},
			entityUri: []string{
				"github.com/AmadlaOrg/Entity",
				"github.com/AmadlaOrg/EntityApplication",
			},
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entityBuild := NewGetService()
			err := entityBuild.Get(&test.paths, test.entityUri)
			if err != nil {
				t.Fatalf("Get failed: %v", err)
			}

			// Perform other assertions and checks as needed

			// Clean up
			err = os.RemoveAll(test.paths.Storage)
			if err != nil {
				t.Fatalf("Cleanup failed: %v", err)
			}
		})
	}
}
