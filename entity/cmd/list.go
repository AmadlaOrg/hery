package cmd

import (
	"fmt"
	"os"

	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/cmd/util"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entities",
	Run: func(cmd *cobra.Command, args []string) {
		entityCmdUtilService := util.NewEntityCmdUtilService()
		err := entityCmdUtilService.Concoct(cmd, args, func(paths *storage.AbsPaths, args []string) {
			entityService := entity.NewEntityService(&gitConfig.Config{})
			entities, err := entityService.CrawlDirectoriesParallel(paths.Entities)
			if err != nil {
				fmt.Println("Error crawling directories:", err)
				return
			}
			displayEntities(entities)
		})
		if err != nil {
			return
		}
	},
}

// displayEntities renders a table in the terminal to easily view a list of the entities
func displayEntities(entities map[string]entity.Entity) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Entity Origin", "Entity Name", "Version")

	for name, e := range entities {
		table.Append(e.Origin, name, e.Version)
	}

	table.Render()
}
