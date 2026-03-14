package database

// New to set up the entity Cache service
func New(dbAbsPath string) Database {
	return &dbImpl{
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
