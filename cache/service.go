package cache

import (
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/cache/parser"
)

// NewCacheService to set up the entity Cache service
//
// To be able to use the cache service it is important to use both `Open` and `Close` methods.
//
// The cache service is using a db connection to handle the storage, so it needs to initialize the connection and
// the closing of the db connection.
func NewCacheService() ICache {
	return &SCache{
		Database: database.NewDatabaseService(),
		Parser:   parser.NewParserService(),
	}
}
