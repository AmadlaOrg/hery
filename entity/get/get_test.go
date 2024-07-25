package get

import (
	"github.com/AmadlaOrg/hery/storage"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		paths     storage.AbsPaths
		entityUri []string
	}{
		{
			name: "Get One",
			paths: storage.AbsPaths{
				Storage:    "/tmp/.hery",
				Collection: "/tmp/.hery/test",
				Entities:   "/tmp/.hery/test/entity",
				Cache:      "/tmp/.hery/test/test.cache",
			},
			entityUri: []string{"github.com/AmadlaOrg/EntityApplication"},
		},
		{
			name: "Get Multiple Paths",
			paths: storage.AbsPaths{
				Storage:    "/tmp/.hery",
				Collection: "/tmp/.hery/test",
				Entities:   "/tmp/.hery/test/entity",
				Cache:      "/tmp/.hery/test/test.cache",
			},
			entityUri: []string{"github.com/AmadlaOrg/Entity", "github.com/AmadlaOrg/EntityApplication"},
		},
		//
		// Invalid
		//
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entityBuild := NewGetService()
			err := entityBuild.Get(&test.paths, test.entityUri)
			if err != nil {
				return
			}
		})
	}
}
