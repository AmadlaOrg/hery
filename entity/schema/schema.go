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

// Schema used by mockery
type Schema interface {
	Load(schemaPath string) (*Definition, error)
	FindSchemaFile(entityPath string) (string, error)
	GenerateSchemaPath(entityPath string) string
	GenerateURN(entityUri string) string

	// Local functions
	loadSchemaFile(schemaPath string) (map[string]any, error)
	mergeSchemas(baseSchema, mainSchema map[string]any) map[string]any
}

// schemaImpl used by mockery
type schemaImpl struct{}

// Help with mocking
var (
	osOpen                = os.Open
	jsonNewDecoder        = json.NewDecoder
	jsonschemaNewCompiler = jsonschema.NewCompiler
	filepathGlob          = filepath.Glob
)

// Load loads the JSON schema from a file and merges it with a base schema
func (s *schemaImpl) Load(schemaPath string) (*Definition, error) {
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
	if idVal, ok := schemaData["$id"].(string); ok {
		schemaId = idVal
	}
	if schemaId == "" {
		log.Printf("Warning: Schema %s has no $id field", schemaName)
	}

	// 5. Return the Schema struct
	return &Definition{
		CompiledSchema: compiledSchema,
		SchemaPath:     schemaPath,
		SchemaName:     schemaName,
		SchemaId:       schemaId,
		Schema:         schemaData,
	}, nil
}

// FindSchemaFile discovers the single *.hery.json schema file in the entity type directory.
// Returns an error if no schema or multiple schemas are found.
func (s *schemaImpl) FindSchemaFile(entityPath string) (string, error) {
	pattern := filepath.Join(entityPath, "*"+EntityJsonSchemaFileExt)
	matches, err := filepathGlob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to search for schema file: %w", err)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no %s schema file found in %s", EntityJsonSchemaFileExt, entityPath)
	}
	if len(matches) > 1 {
		return "", fmt.Errorf("multiple %s schema files found in %s: expected exactly one", EntityJsonSchemaFileExt, entityPath)
	}
	return matches[0], nil
}

// GenerateSchemaPath returns the absolute path for the entity's schema.
// Deprecated: Use FindSchemaFile instead for dynamic schema file discovery.
func (s *schemaImpl) GenerateSchemaPath(entityPath string) string {
	schemaPath, err := s.FindSchemaFile(entityPath)
	if err != nil {
		// Fallback: derive name from directory
		dirName := strings.ToLower(filepath.Base(entityPath))
		return filepath.Join(entityPath, dirName+EntityJsonSchemaFileExt)
	}
	return schemaPath
}

// GenerateURN returns the full URN for a HERY entity type.
// Format: urn:hery:<type-uri-with-/-and-@-replaced-by-:>
func (s *schemaImpl) GenerateURN(entityUri string) string {
	urlToUrn := strings.Replace(entityUri, "/", ":", -1)
	urn := strings.Replace(urlToUrn, "@", ":", 1)
	return fmt.Sprintf("urn:hery:%s", urn)
}

//
// Local functions
//

// loadSchemaFile reads a JSON schema file and returns it as a map
func (s *schemaImpl) loadSchemaFile(schemaPath string) (map[string]any, error) {
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
func (s *schemaImpl) mergeSchemas(baseSchema, mainSchema map[string]any) map[string]any {
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
