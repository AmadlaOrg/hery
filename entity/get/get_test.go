package get

import (
	"fmt"
	"testing"

	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/build"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet_Success(t *testing.T) {
	mockEntity := &entity.MockEntity{}
	mockBuild := &build.MockEntityBuild{}

	paths := &storage.AbsPaths{
		Storage:  "/tmp/hery",
		Entities: "/tmp/hery/entity",
	}

	entityMeta := entity.Entity{
		Uri:     "github.com/AmadlaOrg/TestEntity@v1.0.0",
		Name:    "TestEntity",
		RepoUrl: "https://github.com/AmadlaOrg/TestEntity",
		Origin:  "github.com/AmadlaOrg/TestEntity",
		Version: "v1.0.0",
		AbsPath: "/tmp/hery/entity/github.com/AmadlaOrg/TestEntity@v1.0.0",
		Have:    true,
	}

	mockBuild.EXPECT().Meta(*paths, "github.com/AmadlaOrg/TestEntity@v1.0.0").Return(entityMeta, nil)
	mockEntity.EXPECT().CheckDuplicate(mock.Anything, entityMeta).Return(nil)

	getService := &SGet{
		Entity: mockEntity,
		Build:  mockBuild,
	}

	// Entity has Have=true, so download is skipped (no clone needed)
	err := getService.Get(paths, []string{"github.com/AmadlaOrg/TestEntity@v1.0.0"})
	assert.NoError(t, err)
}

func TestGet_DuplicateError(t *testing.T) {
	mockEntity := &entity.MockEntity{}
	mockBuild := &build.MockEntityBuild{}

	paths := &storage.AbsPaths{
		Storage:  "/tmp/hery",
		Entities: "/tmp/hery/entity",
	}

	entityMeta := entity.Entity{
		Uri:  "github.com/AmadlaOrg/TestEntity@v1.0.0",
		Name: "TestEntity",
	}

	mockBuild.EXPECT().Meta(*paths, "github.com/AmadlaOrg/TestEntity@v1.0.0").Return(entityMeta, nil)
	mockEntity.EXPECT().CheckDuplicate(mock.Anything, entityMeta).Return(fmt.Errorf("duplicate entity found"))

	getService := &SGet{
		Entity: mockEntity,
		Build:  mockBuild,
	}

	err := getService.Get(paths, []string{"github.com/AmadlaOrg/TestEntity@v1.0.0"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")
}

func TestGet_BuildMetaError(t *testing.T) {
	mockBuild := &build.MockEntityBuild{}

	paths := &storage.AbsPaths{
		Storage:  "/tmp/hery",
		Entities: "/tmp/hery/entity",
	}

	mockBuild.EXPECT().Meta(*paths, "invalid-entity").Return(entity.Entity{}, fmt.Errorf("invalid entity URI"))

	getService := &SGet{
		Build: mockBuild,
	}

	err := getService.Get(paths, []string{"invalid-entity"})
	assert.Error(t, err)
}

func TestDownload_AllHave(t *testing.T) {
	getService := &SGet{}

	entities := []entity.Entity{
		{Have: true, Name: "Entity1"},
		{Have: true, Name: "Entity2"},
	}

	paths := &storage.AbsPaths{
		Storage:  "/tmp/hery",
		Entities: "/tmp/hery/entity",
	}

	err := getService.download(paths, entities)
	assert.NoError(t, err)
}
