package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGitService(t *testing.T) {
	t.Run("should return a new instance of Git", func(t *testing.T) {
		gitService := NewGitService()
		assert.NotNil(t, gitService)
		assert.IsType(t, &Git{}, gitService)
	})
}
