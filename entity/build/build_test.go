package build

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuilder_MetaFromRemote(t *testing.T) {
	entityBuildService := NewEntityBuildService()
	remote, err := entityBuildService.MetaFromRemote("amadla", "github.com/AmadlaOrg/Entity")
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.Equal(t, "amadla", remote)
}
