package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntitySchemaService(t *testing.T) {
	t.Run("should return a new instance of Schema", func(t *testing.T) {
		service := New()
		assert.NotNil(t, service)
		assert.IsType(t, &schemaImpl{}, service)
	})
}
