package database

// NewDatabaseService to set up the entity Cache service
func NewDatabaseService() IDatabase {
	return &SDatabase{}
}
