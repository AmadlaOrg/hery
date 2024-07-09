package cmd

import (
	"github.com/AmadlaOrg/hery/client"
	"github.com/spf13/cobra"
)

var ConnectToServerCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to HERY server",
	Run: func(cmd *cobra.Command, args []string) {
		client.Connect()
	},
}
