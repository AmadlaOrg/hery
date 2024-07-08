package cmd

import (
	"github.com/AmadlaOrg/hery/entity"
	"github.com/spf13/cobra"
)

var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query entities",
	Run: func(cmd *cobra.Command, args []string) {
		entity.Query(args)
	},
}
