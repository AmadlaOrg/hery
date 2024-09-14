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
			errMsg:              errors.New("schema validation failed: `id` is empty"),
		},
		{
			name:                "Error: id is invalid format",
			inputId:             "urn:hery:amadla:gith'ub.com:Amadla/Org:Entity:v1.0.0",
			inputCollectionName: "amadla",
			inputEntityUri:      "github.com/AmadlaOrg/Entity@v1.0.0",
			expectedErr:         true,
			errMsg:              errors.New("schema validation failed: invalid `id` format"),
		},
		{
			name:                "Error: collectionName is empty",
			inputId:             "urn:hery:amadla:github.com:AmadlaOrg:Entity:v1.0.0",
			inputCollectionName: "",
			inputEntityUri:      "github.com/AmadlaOrg/Entity@v1.0.0",
			expectedErr:         true,
			errMsg:              errors.New("schema validation failed: invalid `urn` since the other values are not the same as the `_entity` URI: collection:  and entity URI: github.com/AmadlaOrg/Entity@v1.0.0"),
		},
		{
			name:                "Error: collectionName is invalid since it has uppercase",
			inputId:             "urn:hery:amadla:github.com:AmadlaOrg:Entity:v1.0.0",
			inputCollectionName: "Amadla",
			inputEntityUri:      "github.com/AmadlaOrg/Entity@v1.0.0",
			expectedErr:         true,
			errMsg:              errors.New("schema validation failed: invalid `urn` since the other values are not the same as the `_entity` URI: collection: Amadla and entity URI: github.com/AmadlaOrg/Entity@v1.0.0"),
		},
		{
			name:                "Error: entityUri is invalid since it has lowercase where it needs uppercase",
			inputId:             "urn:hery:amadla:github.com:AmadlaOrg:Entity:v1.0.0",
			inputCollectionName: "Amadla",
			inputEntityUri:      "github.com/AmadlaOrg/entity@v1.0.0",
			expectedErr:         true,
			errMsg:              errors.New("schema validation failed: invalid `urn` since the other values are not the same as the `_entity` URI: collection: Amadla and entity URI: github.com/AmadlaOrg/Entity@v1.0.0"),
		},
		{
			name:                "Error: entityUri is invalid since it has lowercase where it needs uppercase",
			inputId:             "urn:hery:amadla:github.com:AmadlaOrg:Entity:v1.0.0",
			inputCollectionName: "Amadla",
			inputEntityUri:      "github.com/AmadlaOrg/entity@v1.0.0",
			expectedErr:         true,
			errMsg:              errors.New("schema validation failed: invalid `urn` since the other values are not the same as the `_entity` URI: collection: Amadla and entity URI: github.com/AmadlaOrg/entity@v1.0.0"),
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
