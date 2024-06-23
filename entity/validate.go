package entity

import "github.com/santhosh-tekuri/jsonschema"

func getSchema() *jsonschema.Schema {
	return nil
}

func entity() error {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7
	return nil
}
