package cmd

import (
	collectionPkgCmd "github.com/AmadlaOrg/hery/collection/cmd"
	entityCmdPkg "github.com/AmadlaOrg/hery/entity/cmd"
	"github.com/spf13/cobra"
)

var EntityCmd = &cobra.Command{
	Use:   "entity",
	Short: "Entity commands",
}

func init() {
	collectionPkgCmd.SetFlags(EntityCmd)
	EntityCmd.AddCommand(entityCmdPkg.ListCmd)
	EntityCmd.AddCommand(entityCmdPkg.GetCmd)
	EntityCmd.AddCommand(entityCmdPkg.ValidateCmd)
}
