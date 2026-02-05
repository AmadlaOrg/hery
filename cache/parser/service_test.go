package parser

import (
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewParserService(t *testing.T) {
	t.Run("should return a new instance of Parser", func(t *testing.T) {
		mockDb := &database.MockCacheDatabase{}
		service := NewParserService(mockDb)
		assert.NotNil(t, service)
		assert.IsType(t, &SParser{}, service)
	})
}
