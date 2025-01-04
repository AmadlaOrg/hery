package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"sync"
)

// IDatabase defines the database interface
type IDatabase interface {
	Initialize(dbPath string) error
	Close() error
	IsInitialized() bool
	CreateTable(table Table)
	Insert(table Table)
	Update(table Table, where []Condition)
	Select(table Table, clauses SelectClauses, joinClauses []JoinClauses)
	Delete(table Table, clauses SelectClauses, joinClauses []JoinClauses)
	DropTable(table Table)
	Apply() error
}

// SDatabase implements IDatabase
type SDatabase struct {
	queries *Queries
}

var (
	db          ISqlDb //*sql.DB
	dbMutex     sync.Locker
	initialized bool
	sqlOpen     = func(driverName, dataSourceName string) (ISqlDb, error) {
		d, err := sql.Open(driverName, dataSourceName)
		return d, err
	}
)

// Initialize establishes the database connection
func (s *SDatabase) Initialize(dbPath string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if db != nil {
		return nil // Already initialized successfully
	}

	var err error
	db, err = sqlOpen("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Set PRAGMA statements for performance
	//
	// Doc: https://stackoverflow.com/questions/57118674/go-sqlite3-with-journal-mode-wal-gives-database-is-locked-error
	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		err := db.Close()
		if err != nil {
			return err
		}
		db = nil
		return fmt.Errorf("error setting journal mode: %v", err)
	}

	// Configure the database connection pool
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// Initialization successful
	initialized = true

	return nil
}

// Close closes the database connection
func (s *SDatabase) Close() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if !initialized {
		// Database was never initialized; nothing to close
		return nil
	}

	if db != nil {
		err := db.Close()
		if err != nil {
			return err
		}
		db = nil
		initialized = false
		return nil
	}

	initialized = false
	return nil
}

// IsInitialized returns true if the database has been initialized
func (s *SDatabase) IsInitialized() bool {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	return initialized
}

func (s *SDatabase) query(
	addTo *[]Query,
	table Table,
	buildQueryFunc func(table Table, columnNames, valuesPlaceholder []string) string,
) {
	for _, row := range table.Rows {
		columnNames, valuesPlaceholder, columnValues := processRow(row)

		// Build the query using the provided function
		query := buildQueryFunc(table, columnNames, valuesPlaceholder)

		// Add the query to the queries list
		s.addQuery(addTo, query, columnValues)
	}
}

// addQuery
func (s *SDatabase) addQuery(slice *[]Query, query string, values []string) {
	*slice = append(*slice, Query{
		Query:  query,
		Values: values,
	})
}

// CreateTable creates a new table
func (s *SDatabase) CreateTable(table Table) {
	var sqlColumns string
	for _, column := range table.Columns {
		var columnConstraints string
		for _, constraint := range column.Constraints {
			columnConstraints = constraint.ToSQL()
		}
		sqlColumn := fmt.Sprintf("\n%s %s %s,", column.ColumnName, column.DataType, columnConstraints)
		sqlColumns = fmt.Sprintf("%s %s", sqlColumns, sqlColumn)
	}

	sqlColumns = strings.TrimSuffix(sqlColumns, " ,")

	var sqlRelationships string
	for _, relationship := range table.Relationships {
		sqlRelationship := fmt.Sprintf(
			",\nFOREIGN KEY(%s) REFERENCES %s(%s)",
			relationship.ColumnName,
			relationship.ReferencesTableName,
			relationship.ReferencesColumnName)
		sqlRelationships = fmt.Sprintf("%s %s", sqlRelationships, sqlRelationship)
	}

	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s%s\n);", table.Name, sqlColumns, sqlRelationships)

	var sqlIndexes string
	for _, index := range table.Columns {
		// TODO: `idx_` what is this? Is this needed?
		//createIndexSQL := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS idx_%s_name ON %s(name);`, table.Name, table.Name)
		sqlIndexe := fmt.Sprintf(
			"CREATE INDEX IF NOT EXISTS idx_%s_%s ON %s(%s);",
			table.Name,
			index.ColumnName,
			table.Name,
			index.ColumnName)
		sqlIndexes = fmt.Sprintf("%s\n%s", sqlIndexes, sqlIndexe)
	}

	queryCreateTable := Query{
		Query: fmt.Sprintf("%s\n%s", createTableSQL, sqlIndexes),
	}

	s.queries.CreateTable = append(s.queries.CreateTable, queryCreateTable)
}

// Insert inserts records into the table
func (s *SDatabase) Insert(table Table) {
	s.query(
		&s.queries.Insert,
		table,
		func(table Table, columnNames, valuesPlaceholder []string) string {
			var b strings.Builder
			b.WriteString("INSERT INTO ")
			b.WriteString(table.Name)
			b.WriteString(" (")
			b.WriteString(strings.Join(columnNames, ", "))
			b.WriteString(") VALUES (")
			b.WriteString(strings.Join(valuesPlaceholder, ", "))
			b.WriteString(");")
			return b.String()
		},
	)
}

// Update updates a record in the table
func (s *SDatabase) Update(table Table, where []Condition) {
	s.query(
		&s.queries.Update,
		table,
		func(table Table, columnNames, valuesPlaceholder []string) string {
			var b strings.Builder
			b.WriteString("UPDATE ")
			b.WriteString(table.Name)
			b.WriteString(" SET ")

			var updates []string
			for _, column := range columnNames {
				updates = append(updates, fmt.Sprintf("%s = ?", column))
			}

			b.WriteString(strings.Join(updates, ", "))
			b.WriteString(buildWhere(where))

			return b.String()
		},
	)
}

// Select retrieves a record from the table
func (s *SDatabase) Select(table Table, clauses SelectClauses, joinClauses []JoinClauses) {
	// Build the SELECT query
	var b strings.Builder
	b.WriteString("SELECT * FROM ")
	b.WriteString(table.Name)

	// Build JOIN clauses, if any
	b.WriteString(" ")
	b.WriteString(buildJoinClauses(joinClauses))

	// Build WHERE, GROUP BY, HAVING, ORDER BY, LIMIT, OFFSET
	b.WriteString(buildWhere(clauses.Where))
	if len(clauses.GroupBy) > 0 {
		// For SELECT, GroupBy is just: GROUP BY col1, col2 ...
		b.WriteString(fmt.Sprintf(" GROUP BY %s", strings.Join(clauses.GroupBy, ", ")))
	}
	b.WriteString(buildHaving(clauses.Having))
	b.WriteString(buildOrderBy(clauses.OrderBy))
	if clauses.Limit != nil {
		b.WriteString(buildLimit(int64(*clauses.Limit)))
	}
	if clauses.Offset != nil {
		b.WriteString(buildOffset(int64(*clauses.Offset)))
	}
	b.WriteString(";")

	// Add the query to s.queries.Select
	s.queries.Select = append(s.queries.Select, Query{
		Query:  b.String(),
		Values: nil, // You could also populate this if you want parameter binding
	})
}

// Delete deletes records from the table
func (s *SDatabase) Delete(table Table, clauses SelectClauses, joinClauses []JoinClauses) {
	var b strings.Builder
	b.WriteString("DELETE FROM ")
	b.WriteString(table.Name)

	// Build JOIN clauses, if you'd like to support them (not standard for SQLite).
	b.WriteString(" ")
	b.WriteString(buildJoinClauses(joinClauses))

	// Build WHERE, GROUP BY, HAVING, ORDER BY, LIMIT, OFFSET
	b.WriteString(buildWhere(clauses.Where))
	// Typically, GROUP BY / HAVING are not valid with DELETE in standard SQL/SQLite,
	// but included here for completeness if your SQL variant allows it.
	if len(clauses.GroupBy) > 0 {
		b.WriteString(fmt.Sprintf(" GROUP BY %s", strings.Join(clauses.GroupBy, ", ")))
	}
	b.WriteString(buildHaving(clauses.Having))
	b.WriteString(buildOrderBy(clauses.OrderBy))
	// LIMIT / OFFSET in DELETE is also non-standard in SQLite,
	// but some other DB engines do allow it.
	if clauses.Limit != nil {
		b.WriteString(buildLimit(int64(*clauses.Limit)))
	}
	if clauses.Offset != nil {
		b.WriteString(buildOffset(int64(*clauses.Offset)))
	}
	b.WriteString(";")

	s.queries.Delete = append(s.queries.Delete, Query{
		Query:  b.String(),
		Values: nil,
	})
}

// DropTable drops the table from the database
func (s *SDatabase) DropTable(table Table) {
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s;", table.Name)
	s.queries.DropTable = append(s.queries.DropTable, Query{
		Query:  dropSQL,
		Values: nil,
	})
}

// Apply executes all queued queries in s.queries
func (s *SDatabase) Apply() error {
	if !s.IsInitialized() {
		return fmt.Errorf(ErrorDatabaseNotInitialized)
	}

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Merge all queries in the desired order.
	var allQueries []string
	for _, q := range s.queries.CreateTable {
		allQueries = append(allQueries, q.Query)
	}
	for _, q := range s.queries.DropTable {
		allQueries = append(allQueries, q.Query)
	}
	for _, q := range s.queries.Insert {
		allQueries = append(allQueries, q.Query)
	}
	for _, q := range s.queries.Update {
		allQueries = append(allQueries, q.Query)
	}
	for _, q := range s.queries.Delete {
		allQueries = append(allQueries, q.Query)
	}
	for _, q := range s.queries.Select {
		allQueries = append(allQueries, q.Query)
	}

	// Execute each query individually within the transaction.
	// This approach is more flexible if you want to handle parameter binding or errors per query.
	for _, queryString := range allQueries {
		_, execErr := tx.Exec(queryString)
		if execErr != nil {
			// Roll back the entire transaction on error
			rbErr := tx.Rollback()
			if rbErr != nil {
				return fmt.Errorf("error rolling back transaction: %w (original error: %v)", rbErr, execErr)
			}
			return fmt.Errorf("error applying query (%s): %w", queryString, execErr)
		}
	}

	// If all queries succeeded, commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	// Optionally, clear the queries after applying
	s.queries = &Queries{}
	return nil
}
