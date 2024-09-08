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
	// 1. Loop all the top properties in the entity.schema.json
	for key, value := range baseSchema {

		// 2. If the property is properties then merge the entity.schema.json base properties in the other one
		if key == "properties" {

			// 2.1: Ensure the "properties" field exists in the main schema
			if _, exists := mainSchema[key]; !exists {
				mainSchema[key] = map[string]any{}
			}

			// 2.2: Merge properties from baseSchema into mainSchema
			baseProperties := value.(map[string]any)
			mainProperties := mainSchema[key].(map[string]any)
			for propertyName, propertyValue := range baseProperties {
				mainProperties[propertyName] = propertyValue
			}
			mainSchema[key] = mainProperties

			// 3. Make sure that the base require properties are in the merge version
		} else if key == "required" {

			// 3.1: Add `require` if it is not in the top properties
			if _, exists := mainSchema[key]; !exists {
				mainSchema[key] = []any{}
			}

			// 3.2: Merge required fields
			baseRequired := value.([]any)
			mainRequired := mainSchema[key].([]any)
			mainRequired = append(mainRequired, baseRequired...)
			mainSchema[key] = mainRequired

			// 4. `additionalProperties` needs to always be set as the same as the base one
		} else if key == "additionalProperties" {
			mainSchema[key] = value
		}
	}

	return mainSchema
}
