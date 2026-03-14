package schema

// New to set up the entity Schema service
func New() Schema {
	return &schemaImpl{}
}
