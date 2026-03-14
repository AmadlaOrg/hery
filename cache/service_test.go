package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCacheService(t *testing.T) {
	t.Run("should return a new instance of Cache", func(t *testing.T) {
		service := New(t.TempDir())
		assert.NotNil(t, service)
		assert.IsType(t, &cacheImpl{}, service)
	})
}
