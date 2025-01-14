package parser

import (
	"fmt"
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/entity/schema"
	"strings"
)

// parseJsonSchemaToSQLiteType
func parseJsonSchemaToSQLiteType(jsonSchemaType schema.DataType) database.DataType {
	switch jsonSchemaType {
	case schema.DataTypeString:
		return database.DataTypeText
	case schema.DataTypeNumber:
		return database.DataTypeNumeric
	case schema.DataTypeInteger:
		return database.DataTypeNumeric
	case schema.DataTypeObject:
		return database.DataTypeText
	case schema.DataTypeArray:
		return database.DataTypeText
	case schema.DataTypeBoolean:
		return database.DataTypeBoolean
	case schema.DataTypeNull:
		return database.DataTypeText // SQLite has no direct equivalent; treat as NULL
	default:
		return database.DataTypeText
	}
}

// parseJsonSchemaFormatToSQLiteType
func parseJsonSchemaFormatToSQLiteType(jsonSchemaDataFormat schema.DataFormat) database.DataType {
	switch jsonSchemaDataFormat {
	case schema.DataFormatDateTime:
		return database.DataTypeDateTime
	case schema.DataFormatTime:
		return database.DataTypeText
	case schema.DataFormatDate:
		return database.DataTypeDate
	case schema.DataFormatDuration:
		return database.DataTypeText
	default:
		return database.DataTypeText
	}
}

// parseJsonSchemaToSQLiteConstraint converts a JSON Schema constraint into an SQLite constraint.
func parseJsonSchemaToSQLiteConstraint(columnName, jsonSchemaConstraint string, value any) string {
	switch jsonSchemaConstraint {
	case "minLength":
		if minLength, ok := value.(int); ok {
			return fmt.Sprintf("LENGTH(%s) >= %d", columnName, minLength)
		}
	case "maxLength":
		if maxLength, ok := value.(int); ok {
			return fmt.Sprintf("LENGTH(%s) <= %d", columnName, maxLength)
		}
	case "enum":
		if values, ok := value.([]any); ok && len(values) > 0 {
			var enumValues []string
			for _, v := range values {
				enumValues = append(enumValues, fmt.Sprintf("'%v'", v))
			}
			return fmt.Sprintf("CHECK (%s IN (%s))", columnName, strings.Join(enumValues, ", "))
		}
	}
	return ""
}

// joinSQLiteConstraints combines multiple SQLite constraints into a single CHECK clause.
func joinSQLiteConstraints(constraints ...string) string {
	if len(constraints) > 0 {
		return fmt.Sprintf("CHECK (%s)", strings.Join(constraints, " AND "))
	}
	return ""
}
