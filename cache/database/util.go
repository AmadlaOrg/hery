package database

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// mergeSqlQueries takes an array of SQL query strings and merges them together
func mergeSqlQueries(sqlQueries *[]string) string {
	return strings.Replace(strings.Join(*sqlQueries, ";\n")+";", ";;", ";", -1)
}

// ToSQL for Column
func (col Column) ToSQL() string {
	var constraints []string
	for _, constraint := range col.Constraints {
		constraints = append(constraints, constraint.ToSQL())
	}

	return fmt.Sprintf("%s %s %s", col.ColumnName, col.DataType, strings.Join(constraints, " "))
}

// ToSQL for Constraint
func (c Constraint) ToSQL() string {
	switch c.Type {
	case ConstraintNotNull:
		return "NOT NULL"
	case ConstraintUnique:
		return "UNIQUE"
	case ConstraintPrimaryKey:
		return "PRIMARY KEY"
	case ConstraintCheck:
		if c.Condition != "" {
			return fmt.Sprintf("CHECK (%s)", c.Condition)
		}
	case ConstraintDefault:
		if c.Default != "" {
			return fmt.Sprintf("DEFAULT %s", c.Default)
		}
	case ConstraintForeignKey:
		if c.References != "" {
			return fmt.Sprintf("FOREIGN KEY REFERENCES %s", c.References)
		}
	case ConstraintAutoincrement:
		return "AUTOINCREMENT"
	}
	return ""
}

// valueToSQL converts Go primitives in a SQL string
func valueToSQL(value any) string {
	switch v := value.(type) {
	case bool:
		// Convert boolean to SQL-compatible value
		if v {
			return "1" // True in SQL
		}
		return "0" // False in SQL
	case string:
		// Wrap strings in single quotes
		return fmt.Sprintf("'%s'", v)
	case int, int64, float64:
		// Convert numeric values directly
		return fmt.Sprintf("%v", v)
	default:
		// Handle NULL or unsupported types
		return "NULL"
	}
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

// Regular expression to validate column name syntax
var validColumnNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// ValidateColumnName validates the column name for SQLite compatibility.
func ValidateColumnName(columnName string) error {
	columnName = strings.TrimSpace(columnName)

	// Check if the column name is empty
	if columnName == "" {
		return errors.New("column name cannot be empty")
	}

	// Check if the column name is a reserved keyword
	if _, isKeyword := sqliteKeywords[strings.ToUpper(columnName)]; isKeyword {
		return fmt.Errorf("column name '%s' is a reserved keyword", columnName)
	}

	// Check if the column name matches the valid regex
	if !validColumnNameRegex.MatchString(columnName) {
		return fmt.Errorf("column name '%s' contains invalid characters", columnName)
	}

	return nil
}
