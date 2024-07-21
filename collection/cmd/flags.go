package cmd

import (
	"errors"
	"github.com/AmadlaOrg/hery/collection/validation"
	"github.com/AmadlaOrg/hery/env"
	"github.com/spf13/cobra"
	"os"
)

var collection string

// A function variable to replace during testing
var validateName = validation.Name

// SetFlags sets the collection flag and uses an env var as a default source
//
// Add the --collection, -c flag to the entity command only if not used in env var
func SetFlags(cmd *cobra.Command) {
	// --collection
	envVarCollection := os.Getenv(env.HeryCollection)
	cmd.PersistentFlags().StringVarP(
		&collection,
		"collection",
		"c",
		envVarCollection,
		"Specify the collection name")
}

// GetCollectionFlag returns the collection name set by the env var or the command flag
func GetCollectionFlag() (string, error) {
	if collection == "" {
		return "", errors.New("no collection name specified (--collection, -c)")
	} else if validateName(collection) {
		return collection, nil
	}
	return "", errors.New("invalid collection")
}
