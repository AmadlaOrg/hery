package schema

import "github.com/santhosh-tekuri/jsonschema/v6"

const (
	EntityJsonSchemaFileName = "schema.hery.json"
	EntityJsonSchemaIdURN    = `^urn:([a-z0-9][a-z0-9-]{0,31}):([a-z0-9][a-z0-9-]+):([a-zA-Z0-9\-.:]+):([a-zA-Z0-9\-.]+)$`
)

type DataType string

// JSON-Schema data type
/*
Convert JSON-Schema Types to SQLite 3 data types:
	string.
	number.
	integer.
	object.
	array.
	boolean.
	null.
*/
const (
	DataTypeString  DataType = "string"
	DataTypeNumber  DataType = "number"
	DataTypeInteger DataType = "integer"
	DataTypeObject  DataType = "object"
	DataTypeArray   DataType = "array"
	DataTypeBoolean DataType = "boolean"
	DataTypeNull    DataType = "null"
)

type DataFormat string

const (
	DataFormatDateTime DataFormat = "date-time"
	DataFormatTime     DataFormat = "time"
	DataFormatDate     DataFormat = "date"
	DataFormatDuration DataFormat = "duration"
)

// Schema different data
type Schema struct {
	CompiledSchema *jsonschema.Schema
	SchemaPath     string
	SchemaName     string
	SchemaId       string
	Schema         map[string]any
}
