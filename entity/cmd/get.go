package cmd

import (
	collectionPkgCmd "github.com/AmadlaOrg/hery/collection/cmd"
	"github.com/AmadlaOrg/hery/entity"
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
		path, err := storage.Path()
		if err != nil {
			return
		}
		entity.Get(collection, path, args)
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
