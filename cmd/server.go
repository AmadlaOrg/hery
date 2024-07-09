package cmd

import (
	serverCmdPkg "github.com/AmadlaOrg/hery/server/cmd"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "HERY Server",
}

func init() {
	ServerCmd.AddCommand(serverCmdPkg.StartServerCmd)
}
