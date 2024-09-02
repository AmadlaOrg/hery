package heryext

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHeryExtService(t *testing.T) {
	t.Run("should return a new instance of HeryExt", func(t *testing.T) {
		service := NewHeryExtService()
		assert.NotNil(t, service)
		assert.IsType(t, &SHeryExt{}, service)
	})
}
