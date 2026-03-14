package validation

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntityValidationService(t *testing.T) {
	t.Run("should return a new instance of Validation", func(t *testing.T) {
		service := New(&gitConfig.Config{})
		assert.NotNil(t, service)
		assert.IsType(t, &validator{}, service)
	})
}
