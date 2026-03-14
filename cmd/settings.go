package cmd

import (
	"github.com/AmadlaOrg/hery/env"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var SettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "List the paths and other environment variables for HERY",
	Run: func(cmd *cobra.Command, args []string) {
		storageService := storage.New()
		heryPath, err := storageService.Main()
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header("Setting", "Value")
		table.Append("Collections path", heryPath)

		envList, err := env.List()
		if err != nil {
			log.Fatal(err)
		}

		for _, varName := range envList {
			val := os.Getenv(varName)
			table.Append(varName, val)
		}

		table.Render()
	},
}
