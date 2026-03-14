package compose

import (
	"testing"

	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/stretchr/testify/assert"
)

func TestNewComposeService(t *testing.T) {
	t.Run("should return a new instance of Composer", func(t *testing.T) {
		service := New(&gitConfig.Config{})
		assert.NotNil(t, service)
		assert.IsType(t, &composer{}, service)
	})
}
