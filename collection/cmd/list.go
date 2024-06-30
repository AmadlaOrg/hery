package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List collections",
	Run: func(cmd *cobra.Command, args []string) {
		/*entityDir, err := storage.Path()
		if err != nil {
			fmt.Println("could not get the root storage directory:", err)
			return
		}

		// TODO: Maybe display that it is downloading
		err = entity.Get(entityDir)
		if err != nil {
			fmt.Println("Error crawling directories:", err)
			return
		}*/
	},
}

// displayCollectionsList display the directories (that are collections) using tablewriter
func displayCollectionsList(collections [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Directories"})

	for _, v := range collections {
		table.Append(v)
	}

	table.Render()
}
