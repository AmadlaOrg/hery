package database

// NewDatabaseService to set up the entity Cache service
func NewDatabaseService() IDatabase {
	return &SDatabase{
		queries: &Queries{
			CreateTable: []Query{},
			DropTable:   []Query{},
			Insert:      []Query{},
			Update:      []Query{},
			Delete:      []Query{},
			Select:      []Query{},
		},
	}
}
