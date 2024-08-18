package cmd

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/get"
	entityValidation "github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
)

var (
	isValidateAll bool
	isRm          bool
)

var ValidateCmd = &cobra.Command{
	Use:   "valid",
	Short: "Validate entity or schemas",
	Run: func(cmd *cobra.Command, args []string) {
		if !isValidateAll && (len(args) == 0 || (len(args) == 0 && isRm)) {
			// Print the usage information (helper message)
			err := cmd.Help()
			if err != nil {
				log.Fatal(err)
			}
			return
		} else if isValidateAll && isRm {
			err := cmd.Help()
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		concoct(cmd, args, func(collectionName string, paths *storage.AbsPaths, args []string) {
			if isRm {
				argsLen := len(args) - 1

				if len(args) == 0 {
					log.Fatal("no entity URI specified")
				} else if argsLen > 60 {
					log.Fatal("too many entity URIs (the limit is 60)")
				}

				//err := TmpEntityCheck(collectionName, args)
				getService := get.NewGetService()
				err := getService.GetInTmp(collectionName, paths, args)
				if err != nil {
					log.Fatal(err)
				}

				println(paths.Entities)
			} else {
				println(args)
				println(isValidateAll)

				for _, arg := range args {
					println(arg)
				}

				entityList, err := entity.CrawlDirectoriesParallel(paths.Entities)
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

// TmpEntityCheck
func TmpEntityCheck(collectionName string, entities []string) error {
	storageService := storage.NewStorageService()

	// Replace paths with temporary directory before .<collectionName>
	tmpPaths, err := storageService.TmpPaths(collectionName)
	if err != nil {
		return err
	}

	err = storageService.MakePaths(*tmpPaths)
	if err != nil {
		return err
	}

	getService := get.NewGetService()
	err = getService.Get(collectionName, tmpPaths, entities)
	if err != nil {
		return errors.New(fmt.Sprintf("error getting entity: %s", err))
	}

	return nil
}
