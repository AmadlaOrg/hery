package cmd

import (
	"os"
	"testing"

	"github.com/AmadlaOrg/hery/env"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestSetFlags(t *testing.T) {
	cmd := &cobra.Command{}
	// Set an environment variable
	err := os.Setenv(env.HeryCollection, "test_env_collection")
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer func(key string) {
		err := os.Unsetenv(key)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
	}(env.HeryCollection)

	SetFlags(cmd)

	flag := cmd.PersistentFlags().Lookup("collection")
	assert.NotNil(t, flag, "--collection flag should be added")
	assert.Equal(t, "test_env_collection", flag.DefValue, "Default value should be set from environment variable")
}

func TestGetCollectionFlag(t *testing.T) {
	tests := []struct {
		name          string
		collectionVal string
		validationVal bool
		expectedVal   string
		expectedErr   string
	}{
		{"No collection specified", "", true, "", "no collection name specified (--collection, -c)"},
		{"Valid collection from flag", "valid_collection", true, "valid_collection", ""},
		{"Invalid collection from flag", "invalid_collection", false, "", "invalid collection"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection = tt.collectionVal

			// Replace validateName with a mock function
			originalValidateName := validateName
			validateName = func(name string) bool {
				return tt.validationVal
			}
			defer func() { validateName = originalValidateName }()

			val, err := GetCollectionFlag()
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedVal, val)
		})
	}
}
