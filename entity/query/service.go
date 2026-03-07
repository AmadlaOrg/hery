package query

import "github.com/AmadlaOrg/hery/cache/database"

// NewQueryService creates a new query service with the given database.
func NewQueryService(db database.IDatabase) IQuery {
	return &SQuery{
		Database: db,
	}
}
