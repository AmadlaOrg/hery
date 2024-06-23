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

func displayEntities(entities map[string]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Entity Name", "Version"})

	for name, version := range entities {
		table.Append([]string{name, version})
	}

	table.Render()
}
