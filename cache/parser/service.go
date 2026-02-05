package parser

import "github.com/AmadlaOrg/hery/cache/database"

// NewParserService to set up the entity Cache service
func NewParserService(db database.IDatabase) IParser {
	return &SParser{
		database: db,
	}
}
