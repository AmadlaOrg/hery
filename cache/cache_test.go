package cache

import (
	"path/filepath"
	"testing"

	"github.com/AmadlaOrg/hery/entity"
	"github.com/stretchr/testify/assert"
)

func TestAddEntity(t *testing.T) {
	cacheService := New(filepath.Join(t.TempDir(), "test.db"))
	err := cacheService.Open()
	assert.NoError(t, err)

	defer func() {
		err := cacheService.Close()
		assert.NoError(t, err)
	}()

	e := &entity.Entity{
		Uri:     "github.com/AmadlaOrg/App/WebServer@v1.0.0",
		Name:    "WebServer",
		RepoUrl: "https://github.com/AmadlaOrg/App",
		Origin:  "github.com/AmadlaOrg/App",
		Version: "v1.0.0",
		AbsPath: "/home/user/.cache/hery/entity/github.com/AmadlaOrg/App/WebServer@v1.0.0",
		Have:    true,
		Exist:   true,
		Content: entity.Content{
			Type: "github.com/AmadlaOrg/App/WebServer@v1.0.0",
			Body: map[string]any{
				"port": 8080,
			},
		},
	}

	err = cacheService.AddEntity(e)
	assert.NoError(t, err)
}

func TestSelectEntity(t *testing.T) {
	cacheService := New(filepath.Join(t.TempDir(), "test.db"))
	err := cacheService.Open()
	assert.NoError(t, err)

	defer func() {
		err := cacheService.Close()
		assert.NoError(t, err)
	}()

	e := &entity.Entity{
		Uri:     "github.com/AmadlaOrg/App/WebServer@v1.0.0",
		Name:    "WebServer",
		RepoUrl: "https://github.com/AmadlaOrg/App",
		Origin:  "github.com/AmadlaOrg/App",
		Version: "v1.0.0",
		AbsPath: "/home/user/.cache/hery/entity/github.com/AmadlaOrg/App/WebServer@v1.0.0",
		Have:    true,
		Exist:   true,
		Content: entity.Content{
			Type: "github.com/AmadlaOrg/App/WebServer@v1.0.0",
			Body: map[string]any{
				"port": 8080,
			},
		},
	}

	err = cacheService.AddEntity(e)
	assert.NoError(t, err)

	result, err := cacheService.SelectEntity("github.com/AmadlaOrg/App/WebServer@v1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, "WebServer", result.Name)
}

func TestSelectEntity_NotFound(t *testing.T) {
	cacheService := New(filepath.Join(t.TempDir(), "test.db"))
	err := cacheService.Open()
	assert.NoError(t, err)

	defer func() {
		err := cacheService.Close()
		assert.NoError(t, err)
	}()

	// Create tables first
	_, err = cacheService.SelectEntity("nonexistent@v1.0.0")
	assert.Error(t, err)
}
