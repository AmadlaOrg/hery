package version

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntityVersionService(t *testing.T) {
	t.Run("should return a new instance of Entity Version", func(t *testing.T) {
		entityVersionService := NewEntityVersionService(&gitConfig.Config{})
		assert.NotNil(t, entityVersionService)
		assert.IsType(t, &SVersion{}, entityVersionService)
	})
}
