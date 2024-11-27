package database

// Errors
const (
	ErrorClosingDatabase        = "closing database"
	ErrorDatabaseNotInitialized = "database not initialized"
)

type DataType string

// Data types
/*
SQLite 3 data type:
	Declared Type	       Type Affinity	Storage Class
	INTEGER	               INTEGER	        INTEGER
	TINYINT, SMALLINT	   INTEGER	        INTEGER
	MEDIUMINT, BIGINT	   INTEGER	        INTEGER
	INT	                   INTEGER	        INTEGER
	REAL, DOUBLE, FLOAT	   REAL	            REAL
	NUMERIC, DECIMAL	   NUMERIC	REAL or INTEGER (if possible)
	TEXT	               TEXT	            TEXT
	CHARACTER, VARCHAR	   TEXT	            TEXT
	CLOB	               TEXT	            TEXT
	BLOB	               BLOB	            BLOB
	BOOLEAN	               NUMERIC	         INTEGER (1 for true, 0 for false)
	DATE, DATETIME	       NUMERIC	TEXT, REAL, or INTEGER depending on the format
*/
const (
	DataTypeInteger    DataType = "INTEGER"
	DataTypeTinyint    DataType = "TINYINT"
	DataTypeBigInteger DataType = "BIGINT"
	DataTypeReal       DataType = "REAL"
	DataTypeNumeric    DataType = "NUMERIC"
	DataTypeDecimal    DataType = "DECIMAL"
	DataTypeBoolean    DataType = "BOOLEAN"
	DataTypeText       DataType = "TEXT"
	DataTypeCharacter  DataType = "CHARACTER"
	DataTypeVarchar    DataType = "VARCHAR"
	DataTypeClob       DataType = "CLOB"
	DataTypeBlob       DataType = "BLOB"
	DataTypeDate       DataType = "DATE"
	DataTypeDateTime   DataType = "DATETIME"
)

// ConstraintType represents a SQLite constraint.
type ConstraintType string

const (
	ConstraintPrimaryKey    ConstraintType = "PRIMARY KEY"
	ConstraintNotNull       ConstraintType = "NOT NULL"
	ConstraintUnique        ConstraintType = "UNIQUE"
	ConstraintCheck         ConstraintType = "CHECK"
	ConstraintDefault       ConstraintType = "DEFAULT"
	ConstraintForeignKey    ConstraintType = "FOREIGN KEY"
	ConstraintAutoincrement ConstraintType = "AUTOINCREMENT"
)

// Table is a basic representation in a struct of a table in a SQL DB
type Table struct {
	Name          string
	Columns       []Column
	Relationships []Relationship
	Rows          []Row
}

// Column is a basic representation in a struct of a column in a SQL DB
type Column struct {
	ColumnName  string
	DataType    DataType
	Constraints []Constraint
}

// Relationship so to create relationships
type Relationship struct {
	ColumnName           string
	ReferencesTableName  string
	ReferencesColumnName string
}

// Row is where the data is being passed compared to the structure in the Column struct
type Row = map[string]any

// Constraint represents a SQLite constraint.
type Constraint struct {
	Type       ConstraintType
	Condition  string // Used for CHECK constraints
	Default    string // Used for DEFAULT values
	References string // Used for FOREIGN KEY references
}

// Queries is used by apply to execute each of the queries in the right order
type Queries struct {
	CreateTable []Query
	DropTable   []Query
	Insert      []Query
	Update      []Query
	Delete      []Query
	Select      []Query
}

// Query contains what is needed when executing a query
type Query struct {
	Query  string
	Values []string
	Result string
}
