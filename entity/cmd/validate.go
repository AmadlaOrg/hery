package cmd

import (
	collectionPkgCmd "github.com/AmadlaOrg/hery/collection/cmd"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/cmd/util"
	"github.com/AmadlaOrg/hery/entity/cmd/validation"
	"github.com/AmadlaOrg/hery/entity/get"
	entityValidation "github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
)

var (
	isValidateAll     bool
	isRm              bool
	getCollectionFlag = collectionPkgCmd.GetCollectionFlag
)

var ValidateCmd = &cobra.Command{
	Use:   "valid",
	Short: "Validate entity or schemas",
	Run: func(cmd *cobra.Command, args []string) {
		if (!isValidateAll && (len(args) == 0 || (len(args) == 0 && isRm))) || isValidateAll && isRm {
			err := cmd.Help()
			if err != nil {
				log.Fatal(err)
			}
			return
		} else if isRm {
			if err := validation.Entities(args); err != nil {
				log.Fatal(err)
			}

			collectionName, err := getCollectionFlag()
			if err != nil {
				log.Fatal(err)
			}

			//err := TmpEntityCheck(collectionName, args)
			getService := get.NewGetService()
			paths, err := getService.GetInTmp(collectionName, args)
			if err != nil {
				log.Fatal(err)
			}

			println(paths.Entities)
			return
		} else if isValidateAll {
			entityCmdUtilService := util.NewEntityCmdUtilService()
			err := entityCmdUtilService.Concoct(cmd, args, func(collectionName string, paths *storage.AbsPaths, args []string) {
				entityService := entity.NewEntityService()
				entityList, err := entityService.CrawlDirectoriesParallel(paths.Entities)
				if err != nil {
					log.Fatal(err)
				}

				if len(entityList) == 0 {
					println("No entity")
				}

				//println(entityList)
				//println(paths.Entities)

				entityValidation := entityValidation.NewEntityValidationService()

				for _, entity := range entityList {
					err := entityValidation.Entity(collectionName, entity.AbsPath)
					if err != nil {
						log.Fatal(err)
						return
					}
					println(entity.AbsPath)
					println(entity.Name)
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
			if err != nil {
				// TODO: Handle error
				return
			}
		}
	},
}

func init() {
	ValidateCmd.PersistentFlags().BoolVarP(
		&isValidateAll,
		"all",
		"a",
		false,
		"Validate all entities")
	ValidateCmd.PersistentFlags().BoolVar(
		&isRm,
		"rm",
		false,
		"Remove entity after validating if it wasn't already downloaded")
}
