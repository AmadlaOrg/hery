package cmd

import (
	collectionPkgCmd "github.com/AmadlaOrg/hery/collection/cmd"
	"github.com/AmadlaOrg/hery/entity/get"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get entity and its dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		collection, err := collectionPkgCmd.GetCollectionFlag()
		if err != nil {
			log.Fatalln(err.Error())
		}
		storageService := storage.NewStorageService()
		path, err := storageService.Main()
		if err != nil {
			return
		}

		// Use the NewGetService function
		getService := get.NewGetService()

		// Call the Get method
		getService.Get(collection, path, args)
		/*if !validation.CollectionName(collectionName) {
			log.Fatalf("Invalid collection name: %s", collectionName)
		}
		storagePath, err := storage.Path()
		if err != nil {
			log.Fatal(err)
		}
		entity.Get("amadla", storagePath, args)*/
	},
}
