package cmd

import (
	clientCmdPkg "github.com/AmadlaOrg/hery/client/cmd"
	"github.com/spf13/cobra"
)

var ClientCmd = &cobra.Command{
	Use:   "client",
	Short: "HERY client",
}

func init() {
	ClientCmd.AddCommand(clientCmdPkg.ConnectToServerCmd)
}
