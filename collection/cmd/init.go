package cmd

import (
	"fmt"
	"github.com/AmadlaOrg/hery/collection/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init collection (create)",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if exactly one argument is provided
		if len(args) != 1 {
			log.Fatal("Missing collection name.")
		}

		arg := args[0]

		// Validate the collection name that is pass in `arg`
		if validation.Name(arg) {
			log.Fatal("Collection name is required or is in the wrong format.")
		}

		// Retrieve storage path
		storageService := storage.NewStorageService()
		storagePath, err := storageService.Main()
		if err != nil {
			log.Fatal(err)
		}

		// Full path to the target directory
		targetDir := fmt.Sprintf("%s/%s", storagePath, arg)

		// Check if the directory exists
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			// Create the directory
			err = os.MkdirAll(fmt.Sprintf("%s/entity", targetDir), os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Directory '%s' with subdirectory 'entity' created.\n", targetDir)
		} else {
			fmt.Printf("Directory '%s' already exists.\n", targetDir)
		}
	},
}
