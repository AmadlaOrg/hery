package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntitySchemaValidationService(t *testing.T) {
	t.Run("should return a new instance of Validation", func(t *testing.T) {
		service := NewEntitySchemaValidationService()
		assert.NotNil(t, service)
		assert.IsType(t, &SValidation{}, service)
	})
}
