package cmd

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entities",
	Run: func(cmd *cobra.Command, args []string) {
		entityDir, err := storage.Path()
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

// displayEntities renders a table in the terminal to easily view a list of the entities
func displayEntities(entities map[string]entity.Entity) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Entity Origin", "Entity Name", "Version"})

	for name, e := range entities {
		table.Append([]string{e.Origin, name, e.Version})
	}

	table.Render()
}
