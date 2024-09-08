package schema

import (
	"encoding/json"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema"
	"log"
	"os"
)

type ISchema interface {
	Load(schemaPath string) (*jsonschema.Schema, error)
	loadSchemaFile(schemaPath string) (map[string]any, error)
	mergeSchemas(baseSchema, mainSchema map[string]any) (map[string]any, error)
}

type SSchema struct{}

var (
	osOpen = os.Open
)

// Load loads the JSON schema from a file and merges it with a base schema
func (s *SSchema) Load(schemaPath string) (*jsonschema.Schema, error) {
	// 1. Load the main schema
	mainSchemaData, err := s.loadSchemaFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load main schema: %w", err)
	}

	// 2. Load the HERY base schema from .schema/entity.schema.json
	baseSchemaPath := ".schema/entity.schema.json"
	baseSchemaData, err := s.loadSchemaFile(baseSchemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load base schema: %w", err)
	}

	// 3. Merge the two schemas
	mergedSchemaData := s.mergeSchemas(baseSchemaData, mainSchemaData)

	println(mergedSchemaData)

	// 4. Compile the merged schema
	/*compiler := jsonschema.NewCompiler()
	mergedSchema, err := compiler.Compile(jsonschema.NewGoLoader(mergedSchemaData))
	if err != nil {
		return nil, fmt.Errorf("failed to compile merged schema: %w", err)
	}*/

	return nil, nil
}

// loadSchemaFile reads a JSON schema file and returns it as a map
func (s *SSchema) loadSchemaFile(schemaPath string) (map[string]any, error) {
	file, err := osOpen(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open schema file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close schema file: %s", err)
		}
	}(file)

	var schemaData map[string]any
	if err = json.NewDecoder(file).Decode(&schemaData); err != nil {
		return nil, fmt.Errorf("failed to decode schema file: %w", err)
	}

	return schemaData, nil
}

// mergeSchemas merges two schemas (base and main) into one
func (s *SSchema) mergeSchemas(baseSchema, mainSchema map[string]any) map[string]any {
	// For simplicity, this example assumes both schemas are maps and merges them at the top level.
	// You may need to handle deeper merging depending on your schema structure.

	// Start with the base schema and add/overwrite properties from the main schema
	for key, value := range mainSchema {
		baseSchema[key] = value
	}

	return baseSchema
}
