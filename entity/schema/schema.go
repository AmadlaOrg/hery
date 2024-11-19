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
	Load(schemaPath string) (*Schema, error)
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
	osOpen                = os.Open
	jsonNewDecoder        = json.NewDecoder
	jsonschemaNewCompiler = jsonschema.NewCompiler
)

// Load loads the JSON schema from a file and merges it with a base schema
func (s *SSchema) Load(schemaPath string) (*Schema, error) {
	// 1. Read the schema file into memory
	schemaData, err := s.loadSchemaFile(schemaPath)
	if err != nil {
		return nil, err
	}

	schemaName := filepath.Base(schemaPath)

	// 2. Create a new compiler
	compiler := jsonschemaNewCompiler()

	// 3. Add the merged schema to the compiler as a resource
	err = compiler.AddResource(schemaName, schemaData) //bytes.NewReader(mergedSchemaJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to add merged schema to the compiler: %w", err)
	}

	// 4. Compile the merged schema
	compiledSchema, err := compiler.Compile(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to compile merged schema: %w", err)
	}

	var schemaId string
	if schemaData["$id"] == nil {
		schemaId = ""
	} else {
		if schemaId = schemaData["$id"].(string); schemaId != "" {
			// TODO: Better handling of warnings
			log.Printf("Warning! Schema $id from %s to %s is empty", schemaName, schemaId)
		}
	}

	// 5. Return the Schema struct
	return &Schema{
		CompiledSchema: compiledSchema,
		SchemaPath:     schemaPath,
		SchemaName:     schemaName,
		SchemaId:       schemaId,
		Schema:         schemaData,
	}, nil
}

// GenerateSchemaPath returns the absolute path for the entity's schema
func (s *SSchema) GenerateSchemaPath(collectionName, entityPath string) string {
	return filepath.Join(entityPath, fmt.Sprintf(".%s", collectionName), EntityJsonSchemaFileName)
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
	// 1. Read the schema file into memory
	file, err := osOpen(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open schema file: %w", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Printf("failed to close schema file: %v", err)
		}
	}(file)

	// 2. Decode the schema
	var schemaData map[string]any
	if err = jsonNewDecoder(file).Decode(&schemaData); err != nil {
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
