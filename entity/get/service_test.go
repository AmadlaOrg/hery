package get

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGetService(t *testing.T) {
	t.Run("should return a new instance of Entity Version", func(t *testing.T) {
		entityVersionService := NewGetService()
		assert.NotNil(t, entityVersionService)
		assert.IsType(t, SGet{}, entityVersionService)
	})
}
