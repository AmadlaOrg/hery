package validation

import (
	"errors"
	"testing"
)

func TestEntities(t *testing.T) {
	tests := []struct {
		name          string
		entities      []string
		expectedError error
	}{
		{
			name:          "No Entities",
			entities:      []string{},
			expectedError: errors.New("no entity URI specified"),
		},
		{
			name:          "Valid Entities",
			entities:      make([]string, 30), // Example of a valid number of entities
			expectedError: nil,
		},
		{
			name:          "Too Many Entities",
			entities:      make([]string, 61),
			expectedError: errors.New("too many entity URIs (the limit is 60)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Entities(tt.entities)
			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
