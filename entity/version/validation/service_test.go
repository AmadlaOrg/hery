package validation

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEntityVersionValidationService(t *testing.T) {
	service := NewEntityVersionValidationService(&gitConfig.Config{})
	assert.NotNil(t, service, "Expected non-nil VersionValidation service")
	//assert.NotNil(t, service.Version, "Expected non-nil Version in VersionValidation")
	assert.IsType(t, &SValidation{}, service, "Expected Version to be of type version.Service")
}
