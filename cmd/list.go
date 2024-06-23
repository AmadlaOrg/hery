package cmd

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entities",
	Run: func(cmd *cobra.Command, args []string) {
		entityDir, err := entity.StorageRoot()
		if err != nil {
			fmt.Println("could not get the root storage directory:", err)
			return
		}

		entities, err := entity.CrawlDirectoriesParallel(entityDir)
		if err != nil {
			fmt.Println("Error crawling directories:", err)
			return
		}

		displayEntities(entities)
	},
}

// displayEntities
func displayEntities(entities map[string]entity.Entity) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Entity Origin", "Entity Name", "Version"})

	for name, e := range entities {
		table.Append([]string{e.Origin, name, e.Version})
	}

	table.Render()
}
