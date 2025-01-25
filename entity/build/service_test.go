package build

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntityBuildService(t *testing.T) {
	t.Run("should return a new instance of Build", func(t *testing.T) {
		service := NewEntityBuildService(&gitConfig.Config{})
		assert.NotNil(t, service)
		assert.IsType(t, &SBuild{}, service)
	})
}
