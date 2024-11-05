package util

import (
	"fmt"
	collectionPkgCmd "github.com/AmadlaOrg/hery/collection/cmd"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
)

type IUtil interface {
	Concoct(
		cmd *cobra.Command,
		args []string,
		handler func(collectionName string, paths *storage.AbsPaths, args []string)) error
}

type SUtil struct {
	newStorageService storage.IStorage
}

// Function variables to allow for easy testing
var (
	getCollectionFlag = collectionPkgCmd.GetCollectionFlag
)

// Concoct sets up the necessary collection and storage paths and executes the provided handler function.
// It retrieves the collection name using the getCollectionFlag function,
// initializes a new storage service, and gets the paths for the specified collection.
// If any errors occur during these steps, they are logged and the handler is not called.
//
// Parameters:
// - cmd: The cobra command that triggered this function.
// - args: The arguments passed to the cobra command.
// - handler: A function that takes the collection name, storage paths (AbsPaths), and arguments, and performs the main logic.
func (s *SUtil) Concoct(
	cmd *cobra.Command,
	args []string,
	handler func(collectionName string, paths *storage.AbsPaths, args []string)) error {
	collectionName, err := getCollectionFlag()
	if err != nil {
		return fmt.Errorf("failed to get collection flag: %w", err)
	}
	paths, err := s.newStorageService.Paths(collectionName)
	if err != nil {
		log.Println("Error getting paths:", err)
		return err
	}
	handler(collectionName, paths, args)
	return nil
}
