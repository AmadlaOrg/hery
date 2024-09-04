package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntityValidationService(t *testing.T) {
	t.Run("should return a new instance of Utils", func(t *testing.T) {
		service := NewEntityCmdUtilService()
		assert.NotNil(t, service)
		assert.IsType(t, &SUtil{}, service)
	})
}
