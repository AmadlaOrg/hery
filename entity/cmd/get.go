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
			if len(args) == 0 {
				log.Fatal("no entity URL specified")
			}

			getService := get.NewGetService()
			err := getService.Get(paths, args)
			if err != nil {
				log.Fatalf("Error getting entity: %s", err)
			}
		})
	},
}
