package cmd

import (
	"github.com/AmadlaOrg/hery/server"
	"github.com/spf13/cobra"
)

var StartServerCmd = &cobra.Command{
	Use:   "start",
	Short: "Start HERY server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Start()
	},
}
