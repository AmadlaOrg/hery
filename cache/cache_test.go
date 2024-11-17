package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddEntity(t *testing.T) {
	cacheService := NewCacheService()
	err := cacheService.Open()
	assert.NoError(t, err)

	defer func(cacheService ICache) {
		err := cacheService.Close()
		if err != nil {

		}
	}(cacheService)

	err = cacheService.AddEntity()
	assert.NoError(t, err)
}
