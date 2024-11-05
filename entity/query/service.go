package query

// NewQueryService to set up the query service
func NewQueryService() IQuery {
	return &SQuery{}
}
