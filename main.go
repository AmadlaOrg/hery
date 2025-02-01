package main

import (
	"github.com/AmadlaOrg/LibraryFramework/cli"
	"github.com/AmadlaOrg/hery/cmd"
	"github.com/spf13/cobra"
)

// main is the entrypoint of this cli application
//
// - It uses cli.New() to setup the default commands and basic flags
// - If there is any commands to add they can be added in the callback function
func main() {
	cli.New(
		"hery",
		"HERY",
		"1.0.0",
		func(rootCmd *cobra.Command) {
			rootCmd.AddCommand(cmd.SettingsCmd)
			rootCmd.AddCommand(cmd.CollectionCmd)
			rootCmd.AddCommand(cmd.ComposeCmd)
			rootCmd.AddCommand(cmd.QueryCmd)
			rootCmd.AddCommand(cmd.EntityCmd)
		})
}
