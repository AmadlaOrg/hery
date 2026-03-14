package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDatabaseService(t *testing.T) {
	t.Run("should return a new instance of Database", func(t *testing.T) {
		service := New("/home/user/")
		assert.NotNil(t, service)
		assert.IsType(t, &dbImpl{}, service)
	})
}
