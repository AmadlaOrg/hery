package version

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntityVersionService(t *testing.T) {
	t.Run("should return a new instance of Entity Version", func(t *testing.T) {
		entityVersionService := NewEntityVersionService()
		assert.NotNil(t, entityVersionService)
		assert.IsType(t, &Service{}, entityVersionService)
	})
}
