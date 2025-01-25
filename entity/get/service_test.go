package get

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGetService(t *testing.T) {
	t.Run("should return a new instance of Entity Version", func(t *testing.T) {
		entityVersionService := NewGetService(&gitConfig.Config{})
		assert.NotNil(t, entityVersionService)
		assert.IsType(t, &SGet{}, entityVersionService)
	})
}
