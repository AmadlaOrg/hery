package remote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGitRemoteService(t *testing.T) {
	t.Run("should return a new instance of GitRemote", func(t *testing.T) {
		gitRemoteService := NewGitRemoteService()
		assert.NotNil(t, gitRemoteService)
		assert.IsType(t, &SRemote{}, gitRemoteService)
	})
}
