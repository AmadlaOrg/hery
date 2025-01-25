package cmd

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
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
		err := entityCmdUtilService.Concoct(cmd, args, func(collectionName string, paths *storage.AbsPaths, args []string) {
			if err := validation.Entities(args); err != nil {
				log.Fatal(err)
			}

			getService := get.NewGetService(&gitConfig.Config{})
			err := getService.Get(collectionName, paths, args)
			if err != nil {
				log.Fatalf("Error getting entity: %s", err)
			}
		})
		if err != nil {
			// TODO: Handle error
			return
		}
	},
}
