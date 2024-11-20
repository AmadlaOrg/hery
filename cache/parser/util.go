package parser

// parseJsonSchemaToSQLiteType
func parseJsonSchemaToSQLiteType(jsonSchemaType string) string {
	var sqliteType string
	switch jsonSchemaType {
	case "string":
		sqliteType = "TEXT"
	case "number":
		sqliteType = "BIGINT"
	case "integer":
		sqliteType = "BIGINT"
	case "object":
		sqliteType = "TEXT"
	case "array":
		sqliteType = "TEXT"
	case "boolean":
		sqliteType = "BOOLEAN"
	case "null":
		sqliteType = "TEXT" // SQLite has no direct equivalent; treat as NULL
	default:
		sqliteType = "TEXT"
	}

	return sqliteType
}

// parseJsonSchemaFormatToSQLiteType
func parseJsonSchemaFormatToSQLiteType(jsonSchemaType string) string {
	var sqliteType string
	switch jsonSchemaType {
	case "date-time":
		sqliteType = "DATETIME"
	case "time":
		sqliteType = "TEXT"
	case "date":
		sqliteType = "DATE"
	case "duration":
		sqliteType = "TEXT"
	}

	return sqliteType
}
