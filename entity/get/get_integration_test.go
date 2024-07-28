package get

import (
	"github.com/AmadlaOrg/hery/storage"
	"os"
	"path/filepath"
	"testing"
)

func Test_Integration_Get(t *testing.T) {
	tests := []struct {
		name       string
		paths      storage.AbsPaths
		entityURIs []string
		collision  bool
		hasError   bool
	}{
		{
			name: "Get One",
			paths: storage.AbsPaths{
				Storage:    filepath.Join(os.TempDir(), ".hery"),
				Collection: filepath.Join(os.TempDir(), ".hery", "collection"),
				Entities:   filepath.Join(os.TempDir(), ".hery", "collection", "testone", "entity"),
				Cache:      filepath.Join(os.TempDir(), ".hery", "collection", "testone", "test.cache"),
			},
			entityURIs: []string{
				"github.com/AmadlaOrg/EntityApplication",
			},
		},
		{
			name: "Get Multiple different URIs",
			paths: storage.AbsPaths{
				Storage:    filepath.Join(os.TempDir(), ".hery"),
				Collection: filepath.Join(os.TempDir(), ".hery", "collection"),
				Entities:   filepath.Join(os.TempDir(), ".hery", "collection", "testmultiple", "entity"),
				Cache:      filepath.Join(os.TempDir(), ".hery", "collection", "testmultiple", "test.cache"),
			},
			entityURIs: []string{
				"github.com/AmadlaOrg/Entity",
				"github.com/AmadlaOrg/EntityApplication",
			},
		},
		{
			name: "Get Multiple identical URIs (pseudo versions)",
			paths: storage.AbsPaths{
				Storage:    filepath.Join(os.TempDir(), ".hery"),
				Collection: filepath.Join(os.TempDir(), ".hery", "collection"),
				Entities:   filepath.Join(os.TempDir(), ".hery", "collection", "testmultiple", "entity"),
				Cache:      filepath.Join(os.TempDir(), ".hery", "collection", "testmultiple", "test.cache"),
			},
			entityURIs: []string{
				"github.com/AmadlaOrg/Entity",
				"github.com/AmadlaOrg/Entity",
			},
		},
		{
			name: "Get Multiple identical URIs (static versions)",
			paths: storage.AbsPaths{
				Storage:    filepath.Join(os.TempDir(), ".hery"),
				Collection: filepath.Join(os.TempDir(), ".hery", "collection"),
				Entities:   filepath.Join(os.TempDir(), ".hery", "collection", "testmultiple", "entity"),
				Cache:      filepath.Join(os.TempDir(), ".hery", "collection", "testmultiple", "test.cache"),
			},
			entityURIs: []string{
				"github.com/AmadlaOrg/Entity@v1.0.0",
				"github.com/AmadlaOrg/Entity@v1.0.0",
			},
		},
		{
			name: "Get Multiple different URIs (static versions)",
			paths: storage.AbsPaths{
				Storage:    filepath.Join(os.TempDir(), ".hery"),
				Collection: filepath.Join(os.TempDir(), ".hery", "collection"),
				Entities:   filepath.Join(os.TempDir(), ".hery", "collection", "testmultiple", "entity"),
				Cache:      filepath.Join(os.TempDir(), ".hery", "collection", "testmultiple", "test.cache"),
			},
			entityURIs: []string{
				"github.com/AmadlaOrg/Entity@v1.0.0",
				"github.com/AmadlaOrg/Entity@v2.0.0",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entityBuild := NewGetService()
			err := entityBuild.Get(&test.paths, test.entityURIs)
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
