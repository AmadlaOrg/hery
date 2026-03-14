package query

import "github.com/AmadlaOrg/hery/cache/database"

// New creates a new query service with the given database.
func New(db database.Database) Query {
	return &queryImpl{
		Database: db,
	}
}
