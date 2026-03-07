package cmd

import (
	"fmt"
	"log"

	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity/compose"
	"github.com/spf13/cobra"
)

var ComposeCmd = &cobra.Command{
	Use:   "compose [entity]@version",
	Short: "Compose the specified entity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entityArg := args[0]
		printToScreen, _ := cmd.Flags().GetBool("print")
		composeService := compose.NewComposeService(&gitConfig.Config{})
		err := composeService.ComposeEntity(entityArg, printToScreen)
		if err != nil {
			fmt.Println("Error:", err)
			log.Fatal(err)
		}
	},
}

func init() {
	ComposeCmd.Flags().BoolP("print", "p", true, "Print composed entity to screen")
}
