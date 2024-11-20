package parser

import (
	"fmt"
	"strings"
)

// parseJsonSchemaToSQLiteType
func parseJsonSchemaToSQLiteType(jsonSchemaType string) string {
	switch jsonSchemaType {
	case "string":
		return "TEXT"
	case "number":
		return "BIGINT"
	case "integer":
		return "BIGINT"
	case "object":
		return "TEXT"
	case "array":
		return "TEXT"
	case "boolean":
		return "BOOLEAN"
	case "null":
		return "TEXT" // SQLite has no direct equivalent; treat as NULL
	default:
		return "TEXT"
	}
}

// parseJsonSchemaFormatToSQLiteType
func parseJsonSchemaFormatToSQLiteType(jsonSchemaType string) string {
	switch jsonSchemaType {
	case "date-time":
		return "DATETIME"
	case "time":
		return "TEXT"
	case "date":
		return "DATE"
	case "duration":
		return "TEXT"
	default:
		return "TEXT"
	}
}

// parseJsonSchemaToSQLiteConstraint
func parseJsonSchemaToSQLiteConstraint(jsonSchemaConstraint, value any) string {
	switch jsonSchemaConstraint {
	case "minLength":
		return fmt.Sprintf("LENGTH(column_name) >= %d", value.(int))
	case "maxLength":
		return fmt.Sprintf("LENGTH(column_name) <= %d", value.(int))
	default:
		return ""
	}
}

// joinSQLiteConstraints
func joinSQLiteConstraints(constraints ...string) string {
	if len(constraints) > 0 {
		return fmt.Sprintf("CHECK (%s)", strings.Join(constraints, " AND "))
	}
	return ""
}
