package cmd

import (
	"github.com/AmadlaOrg/hery/entity/cmd/util"
	"github.com/AmadlaOrg/hery/entity/cmd/validation"
	"github.com/AmadlaOrg/hery/entity/get"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
)

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get entity and its dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		entityCmdUtilService := util.NewEntityCmdUtilService()
		entityCmdUtilService.Concoct(cmd, args, func(collectionName string, paths *storage.AbsPaths, args []string) {
			if err := validation.Entities(args); err != nil {
				log.Fatal(err)
			}

			getService := get.NewGetService()
			err := getService.Get(collectionName, paths, args)
			if err != nil {
				log.Fatalf("Error getting entity: %s", err)
			}
		})
	},
}
