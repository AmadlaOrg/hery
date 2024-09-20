// Package schema contains utils for handling schemas.
package schema

import (
	"encoding/json"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ISchema used by mockery
type ISchema interface {
	Load(schemaPath string) (*jsonschema.Schema, error)
	GenerateSchemaPath(collectionName, entityPath string) string
	GenerateURNPrefix(collectionName string) string
	GenerateURN(urnPrefix, entityUri string) string

	// Local functions
	loadSchemaFile(schemaPath string) (map[string]any, error)
	mergeSchemas(baseSchema, mainSchema map[string]any) map[string]any
}

// SSchema used by mockery
type SSchema struct{}

// Help with mocking
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
	baseSchemaPath, err := filepath.Abs(filepath.Join("..", "..", ".schema", "entity.schema.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve base schema path: %w", err)
	}
	baseSchemaData, err := s.loadSchemaFile(baseSchemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load base schema: %w", err)
	}

	// 3. Merge the two schemas
	mergedSchemaData := s.mergeSchemas(baseSchemaData, mainSchemaData)

	// 4. Create a new compiler
	compiler := jsonschema.NewCompiler()

	// 5. Add the merged schema to the compiler as a resource
	err = compiler.AddResource("merged_schema.json", mergedSchemaData) //bytes.NewReader(mergedSchemaJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to add merged schema to the compiler: %w", err)
	}

	// 6. Compile the merged schema
	schema, err := compiler.Compile("merged_schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile merged schema: %w", err)
	}

	return schema, nil
}

// GenerateSchemaPath returns the absolute path for the entity's schema
func (s *SSchema) GenerateSchemaPath(collectionName, entityPath string) string {
	return filepath.Join(entityPath, "."+collectionName, EntityJsonSchemaFileName)
}

// GenerateURNPrefix returns the URN prefix for JSON-Schema `id`
func (s *SSchema) GenerateURNPrefix(collectionName string) string {
	return fmt.Sprintf("urn:hery:%s:", collectionName)
}

// GenerateURN returns the full URN for JSON-Schema `id`
func (s *SSchema) GenerateURN(urnPrefix, entityUri string) string {
	urlToUrn := strings.Replace(entityUri, "/", ":", 255)
	heryUrnSuffix := strings.Replace(urlToUrn, "@", ":", 1)
	return fmt.Sprintf("%s%s", urnPrefix, heryUrnSuffix)
}

//
// Local functions
//

// loadSchemaFile reads a JSON schema file and returns it as a map
func (s *SSchema) loadSchemaFile(schemaPath string) (map[string]any, error) {
	file, err := osOpen(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open schema file: %w", err)
	}
	defer func(file *os.File) {
		err = file.Close()
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

			// 3. Since all the base properties use `$defs` to globalise the schema definitions
			// they also need to be merged
		} else if key == "$defs" {

			// 3.1: Ensure the "properties" field exists in the main schema
			if _, exists := mainSchema[key]; !exists {
				mainSchema[key] = map[string]any{}
			}

			// 3.2:

			// 4. Make sure that the base require properties are in the merge version
		} else if key == "required" {

			// 4.1: Add `require` if it is not in the top properties
			if _, exists := mainSchema[key]; !exists {
				mainSchema[key] = []any{}
			}

			// 4.2: Merge required fields
			baseRequired := value.([]any)
			mainRequired := mainSchema[key].([]any)
			mainRequired = append(mainRequired, baseRequired...)
			mainSchema[key] = mainRequired

			// 5. `additionalProperties` needs to always be set as the same as the base one
		} else if key == "additionalProperties" {
			mainSchema[key] = value
		}
	}

	// 6. Returns the merged `mainSchema` with the base schema: `.schema/entity.schema.json`
	return mainSchema
}
