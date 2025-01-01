package database

import (
	"errors"
	"fmt"
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

// ValidateOperator validates if the operator is valid.
func ValidateOperator(op string) error {
	if _, valid := validOperators[op]; !valid {
		return fmt.Errorf("invalid operator: %s", op)
	}
	return nil
}

func processRow(row map[string]any) ([]string, []string, []string) {
	var (
		columnNames       []string
		valuesPlaceholder []string
		columnValues      []string
	)

	for key, value := range row {
		columnNames = append(columnNames, key)
		valuesPlaceholder = append(valuesPlaceholder, "?")
		columnValues = append(columnValues, fmt.Sprintf("%v", value))
	}

	return columnNames, valuesPlaceholder, columnValues
}

// buildWhere build the WHERE clause
func buildWhere(where map[string]any) string {
	if len(where) <= 0 {
		return ""
	}
	var whereClauses []string
	for key, value := range where {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = '%v'", key, value))
	}
	return fmt.Sprintf(" WHERE %s", strings.Join(whereClauses, " AND "))
}

// buildGroupBy
func buildGroupBy(orderBy map[string]any) string {
	if len(orderBy) <= 0 {
		return ""
	}
	var orderByClauses []string
	for key, value := range orderBy {
		orderByClauses = append(orderByClauses, fmt.Sprintf("%s %s", key, value))
	}
	return fmt.Sprintf(" GROUP BY %s", strings.Join(orderByClauses, " "))
}

// buildHaving
func buildHaving(having map[string]any) string {
	if len(having) <= 0 {
		return ""
	}
	var havingClauses []string
	for key, value := range having {
		havingClauses = append(havingClauses, fmt.Sprintf("%s = '%v'", key, value))
	}
	return fmt.Sprintf(" HAVING %s", strings.Join(havingClauses, " "))
}

// buildOrderBy
func buildOrderBy(orderBy map[string]any) string {
	if len(orderBy) <= 0 {
		return ""
	}
	var orderByClauses []string
	for key, value := range orderBy {
		orderByClauses = append(orderByClauses, fmt.Sprintf("%s %s", key, value))
	}
	return fmt.Sprintf(" ORDER BY %s", strings.Join(orderByClauses, " "))
}

// buildLimit
func buildLimit(limit int64) string {
	if limit <= 0 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d", limit)
}

// buildOffset
func buildOffset(offset int64) string {
	if offset <= 0 {
		return ""
	}
	return fmt.Sprintf(" OFFSET %d", offset)
}
