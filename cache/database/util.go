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

// buildWhere constructs the WHERE clause from a slice of Condition structs.
func buildWhere(where []Condition) string {
	if len(where) == 0 {
		return ""
	}
	var whereClauses []string
	for _, condition := range where {
		whereClauses = append(whereClauses, fmt.Sprintf("%s %s '%v'", condition.Column, condition.Operator, condition.Value))
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

// buildHaving constructs the HAVING clause from a slice of Condition structs.
func buildHaving(having []Condition) string {
	if len(having) == 0 {
		return ""
	}
	var havingClauses []string
	for _, condition := range having {
		havingClauses = append(havingClauses, fmt.Sprintf("%s %s '%v'", condition.Column, condition.Operator, condition.Value))
	}
	return fmt.Sprintf(" HAVING %s", strings.Join(havingClauses, " AND "))
}

// buildOrderBy constructs the ORDER BY clause from a slice of OrderBy structs.
func buildOrderBy(orderBy []OrderBy) string {
	if len(orderBy) == 0 {
		return ""
	}
	var orderByClauses []string
	for _, ob := range orderBy {
		orderByClauses = append(orderByClauses, fmt.Sprintf("%s %s", ob.Column, ob.Direction))
	}
	return fmt.Sprintf(" ORDER BY %s", strings.Join(orderByClauses, ", "))
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

// buildJoinClauses constructs the JOIN clause(s) from a slice of JoinClauses.
func buildJoinClauses(joins []JoinClauses) string {
	if len(joins) == 0 {
		return ""
	}

	var joinClauses []string
	for _, join := range joins {
		var clause strings.Builder

		// Add join type
		clause.WriteString(string(join.Type))
		clause.WriteString(" JOIN ")

		// Add table name and optional alias
		clause.WriteString(join.Table)
		if join.Alias != nil {
			clause.WriteString(fmt.Sprintf(" AS %s", *join.Alias))
		}

		// Add USING clause if provided
		if len(join.Using) > 0 {
			clause.WriteString(fmt.Sprintf(" USING (%s)", strings.Join(join.Using, ", ")))
		} else if len(join.On) > 0 {
			// Add ON clause if USING is not provided
			var onClauses []string
			for _, condition := range join.On {
				onClauses = append(onClauses, fmt.Sprintf("%s %s %s", condition.Column1, condition.Operator, condition.Column2))
			}
			clause.WriteString(fmt.Sprintf(" ON %s", strings.Join(onClauses, " AND ")))
		}

		// Add additional raw SQL if provided
		if join.Additional != nil {
			clause.WriteString(fmt.Sprintf(" %s", *join.Additional))
		}

		joinClauses = append(joinClauses, clause.String())
	}

	return strings.Join(joinClauses, " ")
}
