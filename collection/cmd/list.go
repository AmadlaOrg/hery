package cmd

import (
	"github.com/AmadlaOrg/hery/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List collections",
	Run: func(cmd *cobra.Command, args []string) {
		storageService := storage.NewStorageService()
		storagePath, err := storageService.Main()
		if err != nil {
			log.Fatal(err)
		}

		// List directories in the storage path
		directories, err := os.ReadDir(storagePath)
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Collections"})

		for _, dir := range directories {
			if dir.IsDir() {
				table.Append([]string{dir.Name()})
			}
		}

		table.Render()
	},
}
