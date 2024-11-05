package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCacheService(t *testing.T) {
	t.Run("should return a new instance of Cache", func(t *testing.T) {
		service := NewCacheService()
		assert.NotNil(t, service)
		assert.IsType(t, &SCache{}, service)
	})
}
