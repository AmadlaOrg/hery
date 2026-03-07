package validation

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestId(t *testing.T) {
	entitySchemaValidationService := NewEntitySchemaValidationService()
	tests := []struct {
		name           string
		inputId        string
		inputEntityUri string
		expectedErr    bool
		errMsg         error
	}{
		{
			name:           "id is valid",
			inputId:        "urn:hery:github.com:AmadlaOrg:Entity:v1.0.0",
			inputEntityUri: "github.com/AmadlaOrg/Entity@v1.0.0",
			expectedErr:    false,
			errMsg:         nil,
		},
		//
		// Error
		//
		{
			name:           "Error: id is empty",
			inputId:        "",
			inputEntityUri: "github.com/AmadlaOrg/Entity@v1.0.0",
			expectedErr:    true,
			errMsg:         errors.New("schema validation failed: `id` is empty"),
		},
		{
			name:           "Error: id is invalid format",
			inputId:        "urn:hery:gith'ub.com:Amadla/Org:Entity:v1.0.0",
			inputEntityUri: "github.com/AmadlaOrg/Entity@v1.0.0",
			expectedErr:    true,
			errMsg:         errors.New("schema validation failed: invalid `id` format"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := entitySchemaValidationService.Id(tt.inputId, tt.inputEntityUri)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
