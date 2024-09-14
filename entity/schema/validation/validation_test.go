package validation

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestId(t *testing.T) {
	entitySchemaValidationService := NewEntitySchemaValidationService()
	tests := []struct {
		name                string
		inputId             string
		inputCollectionName string
		inputEntityUri      string
		expectedErr         bool
		errMsg              error
	}{
		{
			name:                "id is valid",
			inputId:             "urn:hery:amadla:github.com:AmadlaOrg:Entity:v1.0.0",
			inputCollectionName: "amadla",
			inputEntityUri:      "github.com/AmadlaOrg/Entity@v1.0.0",
			expectedErr:         false,
			errMsg:              nil,
		},
		//
		// Error
		//
		{
			name:                "Error: id is empty",
			inputId:             "",
			inputCollectionName: "amadla",
			inputEntityUri:      "github.com/AmadlaOrg/Entity@v1.0.0",
			expectedErr:         true,
			errMsg:              errors.New("schema validation failed: `id` is empty or not set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := entitySchemaValidationService.Id(tt.inputId, tt.inputCollectionName, tt.inputEntityUri)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
