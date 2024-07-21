package cmd

import (
	"github.com/AmadlaOrg/hery/entity/get"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get entity and its dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		concoct(cmd, args, func(collectionName string, paths *storage.AbsPaths, args []string) {
			getService := get.NewGetService()
			err := getService.Get(collectionName, paths, args)
			if err != nil {
				log.Println("Error getting entity:", err)
			}
		})
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
