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

		// Derive schema filename from last segment of entity URI
		// e.g., "amadla.org/entity/application" -> "application.hery.json"
		parts := filepath.Base(entityURI)
		// Strip version if present (e.g., "application@v1.0.0" -> "application")
		if idx := len(parts) - 1; idx >= 0 {
			for i, c := range parts {
				if c == '@' {
					parts = parts[:i]
					break
				}
			}
		}
		schemaFileName := parts + ".hery.json"

		// Create schema file
		schemaPath := filepath.Join(".", schemaFileName)
		schemaContent := fmt.Sprintf(`{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "urn:hery:%s",
  "type": "object",
  "properties": {
    "_type": { "type": "string" },
    "_extends": { "type": "string" },
    "_meta": { "type": "object" },
    "_body": { "type": "object" },
    "_requires": { "type": "array", "items": { "type": "string" } }
  },
  "required": ["_type"]
}
`, entityURI)

		if err := os.WriteFile(schemaPath, []byte(schemaContent), 0644); err != nil {
			log.Fatalf("Failed to create schema file: %v", err)
		}
		fmt.Printf("Created %s\n", schemaPath)

		// Create empty .hery file
		heryPath := filepath.Join(".", "default.hery")
		heryContent := fmt.Sprintf("_type: %s\n_body:\n", entityURI)
		if err := os.WriteFile(heryPath, []byte(heryContent), 0644); err != nil {
			log.Fatalf("Failed to create entity file: %v", err)
		}
		fmt.Printf("Created %s\n", heryPath)
	},
}
