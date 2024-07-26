package build

import (
	"github.com/AmadlaOrg/hery/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Integration_Builder_MetaFromRemote(t *testing.T) {
	paths, err := storage.NewStorageService().Paths("amadla")
	if err != nil {
		t.Fatal(err)
	}

	entityBuildService := NewEntityBuildService()
	remote, err := entityBuildService.MetaFromRemote(*paths, "github.com/AmadlaOrg/EntityApplication")
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, "amadla", remote)
}
