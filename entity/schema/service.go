package schema

// NewEntitySchemaService to set up the entity Schema service
func NewEntitySchemaService() ISchema {
	return &SSchema{}
}
