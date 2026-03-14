package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/entity/query"
	"github.com/spf13/cobra"
)

var QueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query entities",
	RunE: func(cmd *cobra.Command, args []string) error {
		typeFlag, _ := cmd.Flags().GetString("type")
		metaFlag, _ := cmd.Flags().GetString("meta")
		tagFlag, _ := cmd.Flags().GetString("tag")
		jqFlag, _ := cmd.Flags().GetString("jq")

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		dbPath := filepath.Join(homeDir, ".cache", "hery", "hery.db")

		db := database.New(dbPath)
		if err := db.Initialize(); err != nil {
			return fmt.Errorf("failed to open cache database: %w", err)
		}
		defer db.Close()

		queryService := query.New(db)

		results, err := queryService.Query(query.SelectionOpts{
			Type: typeFlag,
			Meta: metaFlag,
			Tag:  tagFlag,
			JQ:   jqFlag,
		})
		if err != nil {
			return err
		}

		output, err := query.FormatJSON(results)
		if err != nil {
			return err
		}
		fmt.Println(output)
		return nil
	},
}

func init() {
	QueryCmd.Flags().String("type", "", "Filter by entity type (glob pattern)")
	QueryCmd.Flags().String("meta", "", "Filter by metadata content (substring)")
	QueryCmd.Flags().String("tag", "", "Filter by tag (substring in meta)")
	QueryCmd.Flags().String("jq", "", "jq expression for transformation")
}
