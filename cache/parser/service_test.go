package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewParserService(t *testing.T) {
	t.Run("should return a new instance of Parser", func(t *testing.T) {
		service := NewParserService()
		assert.NotNil(t, service)
		assert.IsType(t, &SParser{}, service)
	})
}
