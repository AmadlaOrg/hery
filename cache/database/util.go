package database

import (
	"fmt"
	"strings"
)

// mergeSqlQueries takes an array of SQL query strings and merges them together
func mergeSqlQueries(sqlQueries *[]string) string {
	return strings.Replace(strings.Join(*sqlQueries, ";\n")+";", ";;", ";", -1)
}

func (c Constraint) ToSQL(columnName string) string {
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
			return fmt.Sprintf("FOREIGN KEY (%s)", c.References)
		}
	case ConstraintAutoincrement:
		return "AUTOINCREMENT"
	}
	return ""
}
