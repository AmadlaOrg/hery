package query

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewQueryService(t *testing.T) {
	t.Run("should return a new instance of Query", func(t *testing.T) {
		service := NewQueryService()
		assert.NotNil(t, service)
		assert.IsType(t, &SQuery{}, service)
	})
}
