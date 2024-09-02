package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHeryExtService(t *testing.T) {
	t.Run("should return a new instance of Schema", func(t *testing.T) {
		service := NewSchemaService()
		assert.NotNil(t, service)
		assert.IsType(t, &SSchema{}, service)
	})
}
