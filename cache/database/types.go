package database

const (
	ErrorClosingDatabase        = "closing database"
	ErrorDatabaseNotInitialized = "database not initialized"
)

// Table is a basic representation in a struct of a table in a SQL DB
type Table struct {
	Name          string
	Columns       []Column
	Relationships []Relationships
	Rows          []map[string]any
}

// Column is a basic representation in a struct of a column in a SQL DB
type Column struct {
	ColumnName string
	DataType   string
	Constraint string
	Default    string
}

// Relationships so to create relationships
type Relationships struct {
	ColumnName           string
	ReferencesTableName  string
	ReferencesColumnName string
}
