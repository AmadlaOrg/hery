package cmd

import (
	"github.com/AmadlaOrg/hery/entity"
	entityValidation "github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
)

var ValidateCmd = &cobra.Command{
	Use:   "valid",
	Short: "Validate entity or schemas",
	Run: func(cmd *cobra.Command, args []string) {
		concoct(cmd, args, func(collectionName string, paths *storage.AbsPaths, args []string) {
			entityList, err := entity.CrawlDirectoriesParallel(paths.Entities)
			if err != nil {
				log.Fatal(err)
			}

			entityValidation := entityValidation.NewEntityValidationService()

			for _, entity := range entityList {
				err := entityValidation.Entity(collectionName, entity.AbsPath)
				if err != nil {
					log.Fatal(err)
					return
				}
			}

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
