package database

const (
	ErrorClosingDatabase        = "closing database"
	ErrorDatabaseNotInitialized = "database not initialized"
)

// Table is a basic representation in a struct of a table in a SQL DB
type Table struct {
	Name    string
	Columns []Column
	Rows    map[string]string
}

// Column is a basic representation in a struct of a column in a SQL DB
type Column struct {
	ColumnName string
	DataType   string
}
