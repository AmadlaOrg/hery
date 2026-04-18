package cmd

import (
	"fmt"
	"log"
	"os"

	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/AmadlaOrg/hery/entity/compose"
	"github.com/AmadlaOrg/hery/entity/resolve"
	"github.com/spf13/cobra"
)

var ComposeCmd = &cobra.Command{
	Use:   "compose [entity]@version",
	Short: "Compose the specified entity",
	Long:  "Compose the specified entity. Use --dir to compose a local directory of .hery files instead of a cached entity URI.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("dir")
		layer, _ := cmd.Flags().GetInt("layer")

		if dir != "" {
			return runDirCompose(dir, layer)
		}

		if len(args) == 0 {
			return fmt.Errorf("entity argument required when --dir is not set")
		}
		printToScreen, _ := cmd.Flags().GetBool("print")
		composeService := compose.New(&gitConfig.Config{})
		if err := composeService.ComposeEntity(args[0], printToScreen); err != nil {
			log.Fatal(err)
		}
		return nil
	},
}

func runDirCompose(dir string, layer int) error {
	r := resolve.New()
	result, err := r.Resolve(dir)
	if err != nil {
		return err
	}

	var out []byte
	if layer > 0 {
		if layer > len(result.Layers) {
			return fmt.Errorf("--layer %d requested but only %d layers resolved", layer, len(result.Layers))
		}
		out, err = resolve.Marshal(result.Layers[layer-1].Docs)
	} else {
		out, err = resolve.MarshalAll(result.Layers)
	}
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(out)
	return err
}

func init() {
	ComposeCmd.Flags().BoolP("print", "p", true, "Print composed entity to screen")
	ComposeCmd.Flags().String("dir", "", "Compose a local directory of .hery files")
	ComposeCmd.Flags().Int("layer", 0, "When set with --dir, emit only the Nth merge layer (1-indexed)")
}
