package cmd

import (
	collectionPkgCmd "github.com/AmadlaOrg/hery/collection/cmd"
	"github.com/spf13/cobra"
)

var CollectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "Collections",
}

func init() {
	CollectionCmd.AddCommand(collectionPkgCmd.ListCmd)
	CollectionCmd.AddCommand(collectionPkgCmd.InitCmd)
}
