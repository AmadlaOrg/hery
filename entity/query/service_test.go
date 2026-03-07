package query

import (
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewQueryService(t *testing.T) {
	t.Run("should return a new instance of Query", func(t *testing.T) {
		mockDb := &database.MockCacheDatabase{}
		service := NewQueryService(mockDb)
		assert.NotNil(t, service)
		assert.IsType(t, &SQuery{}, service)
	})
}
