package cmd

import (
	"errors"
	"github.com/AmadlaOrg/hery/collection/validation"
	"github.com/AmadlaOrg/hery/env"
	"github.com/spf13/cobra"
	"os"
)

var collection string

// SetFlags sets the collection flag and uses an env var as a default source
//
// Add the --collection flag to the entity command only if not used in env var
func SetFlags(cmd *cobra.Command) {
	// --collection
	envVarCollection := os.Getenv(env.HeryCollection)
	cmd.PersistentFlags().StringVarP(&collection, "collection", "c", envVarCollection, "Specify the collection")
}

// GetCollectionFlag returns the collection name set by the env var or the command flag
func GetCollectionFlag() (string, error) {
	if validation.CollectionName(collection) {
		return collection, nil
	}
	return "", errors.New("invalid collection")
}
