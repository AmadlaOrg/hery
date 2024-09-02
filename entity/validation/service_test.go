package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntityValidationService(t *testing.T) {
	t.Run("should return a new instance of Validation", func(t *testing.T) {
		service := NewEntityValidationService()
		assert.NotNil(t, service)
		assert.IsType(t, &SValidation{}, service)
	})
}
