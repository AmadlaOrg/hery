package compose

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestComposeEntity_PrintToScreen(t *testing.T) {
	mockStorage := &storage.MockStorage{}
	mockEntity := &entity.MockEntity{}

	paths := &storage.AbsPaths{
		Storage:  "/tmp/hery",
		Entities: "/tmp/hery/entity",
	}

	mockStorage.EXPECT().Paths().Return(paths, nil)
	mockEntity.EXPECT().ReadAll("/tmp/hery/entity").Return([]map[string]any{
		{
			"_type": "example.com/App@v1.0.0",
			"_body": map[string]any{
				"name": "test-app",
				"port": 8080,
			},
		},
	}, nil)

	composer := &composer{
		Storage: mockStorage,
		Entity:  mockEntity,
	}

	err := composer.ComposeEntity("example.com/App@v1.0.0", true)
	assert.NoError(t, err)
}

func TestComposeEntity_WriteToFile(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	mockStorage := &storage.MockStorage{}
	mockEntity := &entity.MockEntity{}

	paths := &storage.AbsPaths{
		Storage:  "/tmp/hery",
		Entities: "/tmp/hery/entity",
	}

	mockStorage.EXPECT().Paths().Return(paths, nil)
	mockEntity.EXPECT().ReadAll("/tmp/hery/entity").Return([]map[string]any{
		{
			"_type": "example.com/App@v1.0.0",
			"_body": map[string]any{
				"name": "test-app",
			},
		},
	}, nil)

	composer := &composer{
		Storage: mockStorage,
		Entity:  mockEntity,
	}

	err := composer.ComposeEntity("example.com/App@v1.0.0", false)
	assert.NoError(t, err)

	// Verify file was created
	_, statErr := os.Stat(filepath.Join(tmpDir, "composed.hery"))
	assert.NoError(t, statErr)
}

func TestComposeEntity_WithExtendsResolution(t *testing.T) {
	mockStorage := &storage.MockStorage{}
	mockEntity := &entity.MockEntity{}

	paths := &storage.AbsPaths{
		Storage:  "/tmp/hery",
		Entities: "/tmp/hery/entity",
	}

	mockStorage.EXPECT().Paths().Return(paths, nil)
	// First ReadAll: reads from the entities root directory
	mockEntity.EXPECT().ReadAll("/tmp/hery/entity").Return([]map[string]any{
		{
			"_type":   "example.com/App@v1.0.0",
			"_extends": "github.com/some-org/base-configs/app",
			"_body": map[string]any{
				"port": 9090,
			},
		},
	}, nil)
	// Second ReadAll: resolves _extends by reading from the extended entity's cached directory
	mockEntity.EXPECT().ReadAll("/tmp/hery/entity/github.com/some-org/base-configs/app").Return([]map[string]any{
		{
			"_type": "example.com/App@v1.0.0",
			"_body": map[string]any{
				"color": "blue",
				"port":  8080,
			},
		},
	}, nil)

	composer := &composer{
		Storage: mockStorage,
		Entity:  mockEntity,
	}

	err := composer.ComposeEntity("example.com/App@v1.0.0", true)
	assert.NoError(t, err)
}

func TestComposeEntity_NoEntities(t *testing.T) {
	mockStorage := &storage.MockStorage{}
	mockEntity := &entity.MockEntity{}

	paths := &storage.AbsPaths{
		Storage:  "/tmp/hery",
		Entities: "/tmp/hery/entity",
	}

	mockStorage.EXPECT().Paths().Return(paths, nil)
	mockEntity.EXPECT().ReadAll(mock.Anything).Return([]map[string]any{}, nil)

	composer := &composer{
		Storage: mockStorage,
		Entity:  mockEntity,
	}

	err := composer.ComposeEntity("anything", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no entities found")
}

func TestComposeEntity_StorageError(t *testing.T) {
	mockStorage := &storage.MockStorage{}
	mockStorage.EXPECT().Paths().Return(nil, assert.AnError)

	composer := &composer{
		Storage: mockStorage,
	}

	err := composer.ComposeEntity("anything", true)
	assert.Error(t, err)
}
