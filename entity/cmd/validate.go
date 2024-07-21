package cmd

import (
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:   "valid",
	Short: "Validate entity or schemas",
	Run: func(cmd *cobra.Command, args []string) {
		concoct(cmd, args, func(collectionName string, paths *storage.AbsPaths, args []string) {
			// Add your validation logic here
			// entityDir, err := storage.Path()
			// if err != nil {
			//     fmt.Println("could not get the root storage directory:", err)
			//     return
			// }
			// err = entity.Validate(entityDir)
			// if err != nil {
			//     fmt.Println("Error validating entities:", err)
			//     return
			// }
		})
	},
}
