package database

// NewDatabaseService to set up the entity Cache service
func NewDatabaseService(dbAbsPath string) IDatabase {
	return &SDatabase{
		dbAbsPath: dbAbsPath,
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
