package entity

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntityService(t *testing.T) {
	t.Run("should return a new instance of Builder", func(t *testing.T) {
		service := NewEntityService(&gitConfig.Config{})
		assert.NotNil(t, service)
		assert.IsType(t, &SEntity{}, service)
	})
}
