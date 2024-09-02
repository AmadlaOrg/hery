package schema

import (
	"fmt"
	"github.com/santhosh-tekuri/jsonschema"
	"os"
)

type ISchema interface {
	Load(schemaPath string) (*jsonschema.Schema, error)
}

type SSchema struct{}

// Load loads the JSON schema from a file
func (s *SSchema) Load(schemaPath string) (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()
	openSchemaPath, err := os.Open(schemaPath)
	if err != nil {
		return nil, err
	}
	if err := compiler.AddResource("schema.json", openSchemaPath); err != nil {
		return nil, fmt.Errorf("failed to load schema: %w", err)
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}
	return schema, nil
}
