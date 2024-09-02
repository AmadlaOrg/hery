package compose

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntityValidationService(t *testing.T) {
	t.Run("should return a new instance of Composer", func(t *testing.T) {
		service := NewComposeService()
		assert.NotNil(t, service)
		assert.IsType(t, &Composer{}, service)
	})
}
