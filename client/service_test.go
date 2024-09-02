package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClientService(t *testing.T) {
	t.Run("should return a new instance of Client", func(t *testing.T) {
		service := NewClientService()
		assert.NotNil(t, service)
		assert.IsType(t, &SClient{}, service)
	})
}
