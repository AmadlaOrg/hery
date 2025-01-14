package database

import "regexp"

// Errors
const (
	ErrorClosingDatabase        = "closing database"
	ErrorDatabaseNotInitialized = "database not initialized"
)

type Operator string

const (
	OperatorEqual              Operator = "="
	OperatorIntegerNotEqual    Operator = "<>"
	OperatorNotEqual           Operator = "!="
	OperatorGreaterThan        Operator = ">"
	OperatorGreaterThanOrEqual Operator = ">="
	OperatorLessThan           Operator = "<"
	OperatorLessThanOrEqual    Operator = "<="
	OperatorLike               Operator = "LIKE"
	OperatorILike              Operator = "ILIKE"
	OperatorRLike              Operator = "RLIKE"
	OperatorNotLike            Operator = "NOT LIKE"
	OperatorNotILike           Operator = "NOT ILIKE"
	OperatorNotRLike           Operator = "NOT RLIKE"
	OperatorGlob               Operator = "GLOB"
	OperatorIs                 Operator = "IS"
	OperatorNotIs              Operator = "NOT IS"
	OperatorIn                 Operator = "IN"
	OperatorNotIn              Operator = "NOT IN"
	OperatorBetween            Operator = "BETWEEN"
	OperatorNotBetween         Operator = "NOT BETWEEN"
	OperatorExists             Operator = "EXISTS"
	OperatorMatch              Operator = "MATCH"
	OperatorRegexp             Operator = "REGEXP"
)

// ValidOperators is a map of allowed operators for fast lookup.
var validOperators = map[string]struct{}{
	"=":           {},
	"<>":          {},
	"!=":          {},
	"<":           {},
	"<=":          {},
	">":           {},
	">=":          {},
	"LIKE":        {},
	"ILIKE":       {},
	"RLIKE":       {},
	"NOT LIKE":    {},
	"NOT ILIKE":   {},
	"NOT RLIKE":   {},
	"GLOB":        {},
	"IS":          {},
	"IS NOT":      {},
	"IN":          {},
	"NOT IN":      {},
	"BETWEEN":     {},
	"NOT BETWEEN": {},
	"EXISTS":      {},
	"MATCH":       {},
	"REGEXP":      {},
}

// List of SQLite reserved keywords
var sqliteKeywords = map[string]struct{}{
	"ABORT": {}, "ACTION": {}, "ADD": {}, "AFTER": {}, "ALL": {},
	"ALTER": {}, "ANALYZE": {}, "AND": {}, "AS": {}, "ASC": {},
	"ATTACH": {}, "AUTOINCREMENT": {}, "BEFORE": {}, "BEGIN": {},
	"BETWEEN": {}, "BY": {}, "CASCADE": {}, "CASE": {}, "CAST": {},
	"CHECK": {}, "COLLATE": {}, "COLUMN": {}, "COMMIT": {}, "CONFLICT": {},
	"CONSTRAINT": {}, "CREATE": {}, "CROSS": {}, "CURRENT_DATE": {},
	"CURRENT_TIME": {}, "CURRENT_TIMESTAMP": {}, "DATABASE": {}, "DEFAULT": {},
	"DEFERRABLE": {}, "DEFERRED": {}, "DELETE": {}, "DESC": {}, "DETACH": {},
	"DISTINCT": {}, "DROP": {}, "EACH": {}, "ELSE": {}, "END": {},
	"ESCAPE": {}, "EXCEPT": {}, "EXCLUSIVE": {}, "EXISTS": {}, "EXPLAIN": {},
	"FAIL": {}, "FOR": {}, "FOREIGN": {}, "FROM": {}, "FULL": {}, "GLOB": {},
	"GROUP": {}, "HAVING": {}, "IF": {}, "IGNORE": {}, "IMMEDIATE": {},
	"IN": {}, "INDEX": {}, "INDEXED": {}, "INITIALLY": {}, "INNER": {},
	"INSERT": {}, "INSTEAD": {}, "INTERSECT": {}, "INTO": {}, "IS": {},
	"ISNULL": {}, "JOIN": {}, "KEY": {}, "LEFT": {}, "LIKE": {}, "LIMIT": {},
	"MATCH": {}, "NATURAL": {}, "NO": {}, "NOT": {}, "NOTNULL": {}, "NULL": {},
	"OF": {}, "OFFSET": {}, "ON": {}, "OR": {}, "ORDER": {}, "OUTER": {},
	"PLAN": {}, "PRAGMA": {}, "PRIMARY": {}, "QUERY": {}, "RAISE": {},
	"REFERENCES": {}, "REGEXP": {}, "REINDEX": {}, "RELEASE": {}, "RENAME": {},
	"REPLACE": {}, "RESTRICT": {}, "RIGHT": {}, "ROLLBACK": {}, "ROW": {},
	"SAVEPOINT": {}, "SELECT": {}, "SET": {}, "TABLE": {}, "TEMP": {},
	"TEMPORARY": {}, "THEN": {}, "TO": {}, "TRANSACTION": {}, "TRIGGER": {},
	"UNION": {}, "UNIQUE": {}, "UPDATE": {}, "USING": {}, "VACUUM": {},
	"VALUES": {}, "VIEW": {}, "VIRTUAL": {}, "WHEN": {}, "WHERE": {}, "WITH": {},
	"WITHOUT": {},
}

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

// Regular expression to validate column name syntax
var validColumnNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

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

type SelectClauses struct {
	Where   []Condition
	GroupBy []string
	Having  []Condition
	OrderBy []OrderBy
	Limit   *int
	Offset  *int
}

type Condition struct {
	Column   string
	Operator Operator
	Value    any
}

type OrderByDirection string

const (
	OrderByAsc  OrderByDirection = "ASC"
	OrderByDesc OrderByDirection = "DESC"
)

type OrderBy struct {
	Column    string
	Direction OrderByDirection
}

type JoinClauses struct {
	Table      string          // Name of the table to join
	Alias      *string         // Optional alias for the table
	Type       JoinType        // Type of join: INNER, LEFT, RIGHT, etc.
	On         []JoinCondition // Conditions for the ON clause
	Using      []string        // Columns to use for the USING clause (alternative to ON)
	Additional *string         // Optional additional raw SQL for advanced use
}

// JoinType enumerates the types of JOIN operations
type JoinType string

const (
	JoinTypeInner JoinType = "INNER"
	JoinTypeLeft  JoinType = "LEFT"
	JoinTypeRight JoinType = "RIGHT"
	JoinTypeFull  JoinType = "FULL"
	JoinTypeCross JoinType = "CROSS"
)

// JoinCondition struct for expressing conditions in the ON clause
type JoinCondition struct {
	Column1  string   // First column in the condition (with optional table prefix)
	Operator Operator // Operator for the condition (e.g., '=', '<', etc.)
	Column2  string   // Second column in the condition (with optional table prefix)
}

/*
Common SQL Clauses

    SELECT Clauses:
        WHERE: Filters rows based on conditions.
        GROUP BY: Groups rows sharing a property for aggregation.
        HAVING: Filters groups based on aggregate conditions.
        ORDER BY: Sorts the result set.
        LIMIT (or FETCH/TOP in some databases): Limits the number of rows.
        OFFSET: Skips rows in the result set.

    Join Clauses:
        INNER JOIN: Combines rows from two tables where a condition is met.
        LEFT JOIN (or LEFT OUTER JOIN): Includes all rows from the left table and matched rows from the right.
        RIGHT JOIN (or RIGHT OUTER JOIN): Includes all rows from the right table and matched rows from the left.
        FULL JOIN (or FULL OUTER JOIN): Includes rows from both tables, with NULL for non-matches.

    Subquery Clauses:
        EXISTS: Checks if a subquery returns any rows.
        IN: Filters rows where a value matches any in a list or subquery.

    Other Clauses:
        UNION / UNION ALL: Combines result sets from two queries.
        EXCEPT: Returns rows from the first query not in the second.
        INTERSECT: Returns rows common to both queries.
        WITH (Common Table Expression or CTE): Defines a temporary named result set.
*/

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
