package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init [entity-uri]",
	Short: "Initialize a new entity",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entityURI := args[0]

		// Create schema file
		schemaPath := filepath.Join(".", "schema.hery.json")
		schemaContent := fmt.Sprintf(`{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "urn:hery:%s",
  "type": "object",
  "properties": {
    "_type": { "type": "string" },
    "_extends": { "type": "string" },
    "_meta": { "type": "object" },
    "_body": { "type": "object" }
  },
  "required": ["_type"]
}
`, entityURI)

		if err := os.WriteFile(schemaPath, []byte(schemaContent), 0644); err != nil {
			log.Fatalf("Failed to create schema file: %v", err)
		}
		fmt.Printf("Created %s\n", schemaPath)

		// Create empty .hery file
		heryPath := filepath.Join(".", "entity.hery")
		heryContent := fmt.Sprintf("_type: %s\n_body:\n", entityURI)
		if err := os.WriteFile(heryPath, []byte(heryContent), 0644); err != nil {
			log.Fatalf("Failed to create entity file: %v", err)
		}
		fmt.Printf("Created %s\n", heryPath)
	},
}
