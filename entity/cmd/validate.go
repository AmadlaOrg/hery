package cmd

import (
	"fmt"
	"log"

	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	entityPkg "github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/cmd/util"
	"github.com/AmadlaOrg/hery/entity/cmd/validation"
	"github.com/AmadlaOrg/hery/entity/get"
	entityValidation "github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
)

var (
	isValidateAll bool
	isRm          bool
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

			getService := get.NewGetService(&gitConfig.Config{})
			paths, err := getService.GetInTmp(args)
			if err != nil {
				log.Fatal(err)
			}

			println(paths.Entities)
			return
		} else if isValidateAll {
			gitCfg := &gitConfig.Config{}
			entityCmdUtilService := util.NewEntityCmdUtilService()
			err := entityCmdUtilService.Concoct(cmd, args, func(paths *storage.AbsPaths, args []string) {
				entityService := entityPkg.NewEntityService(gitCfg)
				entityList, err := entityService.CrawlDirectoriesParallel(paths.Entities)
				if err != nil {
					log.Fatal(err)
				}

				if len(entityList) == 0 {
					fmt.Println("No entities found")
					return
				}

				validationService := entityValidation.NewEntityValidationService(gitCfg)
				for name, e := range entityList {
					docs, readErr := entityService.ReadAll(e.AbsPath)
					if readErr != nil {
						fmt.Printf("Error reading entity %s: %v\n", name, readErr)
						continue
					}
					for _, doc := range docs {
						if valErr := validationService.Entity(nil, doc); valErr != nil {
							fmt.Printf("Validation failed for %s: %v\n", name, valErr)
						} else {
							fmt.Printf("Entity %s: valid\n", name)
						}
					}
				}
			})
			if err != nil {
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
