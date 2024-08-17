package cmd

import (
	"github.com/AmadlaOrg/hery/collection/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init collection (create)",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if exactly one argument is provided
		if len(args) != 1 {
			log.Fatal("Missing collection name.")
		}

		collectionName := args[0]

		// Validate the collection name that is pass in `arg`
		if validation.Name(collectionName) {
			log.Fatal("Collection name is required or is in the wrong format.")
		}

		// Retrieve storage path
		storageService := storage.NewStorageService()
		paths, err := storageService.Paths(collectionName)
		if err != nil {
			log.Fatal(err)
		}

		err = storageService.MakePaths(*paths)
		if err != nil {
			return
		}

		log.Printf("Collection '%s' created within the storage: '%s'.", collectionName, paths.Storage)
	},
}
