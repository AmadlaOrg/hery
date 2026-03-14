package parser

import "github.com/AmadlaOrg/hery/cache/database"

// New to set up the entity Cache service
func New(db database.Database) Parser {
	return &parser{
		database: db,
	}
}
