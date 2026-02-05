package collection

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCollectionService(t *testing.T) {
	t.Run("should return a new instance of Collection", func(t *testing.T) {
		service := NewCollectionService(&gitConfig.Config{})
		assert.NotNil(t, service)
		assert.IsType(t, &SCollection{}, service)
	})
}
