package build

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntityBuildService(t *testing.T) {
	t.Run("should return a new instance of Builder", func(t *testing.T) {
		service := NewEntityBuildService()
		assert.NotNil(t, service)
		assert.IsType(t, &Builder{}, service)
	})
}
