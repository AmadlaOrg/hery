package cmd

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity/compose"
	"github.com/spf13/cobra"
	"os"
)

var ComposeCmd = &cobra.Command{
	Use:   "compose [entity]@version",
	Short: "Compose the specified entity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entityArg := args[0]
		printToScreen, _ := cmd.Flags().GetBool("print")
		composeService := compose.NewComposeService()
		err := composeService.ComposeEntity(entityArg, printToScreen)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}
