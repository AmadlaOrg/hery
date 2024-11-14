package collection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCollectionService(t *testing.T) {
	t.Run("should return a new instance of Collection", func(t *testing.T) {
		service := NewCollectionService()
		assert.NotNil(t, service)
		assert.IsType(t, &SCollection{}, service)
	})
}
