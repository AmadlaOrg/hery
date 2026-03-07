package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
	"sync"
)

// IDatabase defines the database interface
type IDatabase interface {
	Initialize() error
	Close() error
	IsInitialized() bool
	CreateTable()
	Insert(table Table)
	Update(table Table, where []Condition)
	Select(table Table, clauses SelectClauses, joinClauses []JoinClauses)
	Delete(table Table, clauses SelectClauses)
	DeleteDb() error
	Apply() error
	QueryRows(query string, args ...any) ([]map[string]any, error)
}

// SDatabase implements IDatabase
type SDatabase struct {
	dbAbsPath   string
	queries     *Queries
	sqlDB       ISqlDb
	initialized bool
}

var (
	db          ISqlDb
	dbMutex     sync.Mutex // sync.Locker
	initialized bool
	sqlOpen     = func(driverName, dataSourceName string) (ISqlDb, error) {
		return sql.Open(driverName, dataSourceName)
	}
	osRemove = os.Remove
)

//go:embed table/table.sql
var sqlTables string

// Initialize establishes the database connection
func (s *SDatabase) Initialize() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if db != nil {
		return nil // Already initialized successfully
	}

	var err error
	db, err = sqlOpen("sqlite3", s.dbAbsPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Set PRAGMA statements for performance
	//
	// Doc: https://stackoverflow.com/questions/57118674/go-sqlite3-with-journal-mode-wal-gives-database-is-locked-error
	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		pragmaErr := err
		if closeErr := db.Close(); closeErr != nil {
			return fmt.Errorf("error setting journal mode: %v (also failed to close: %v)", pragmaErr, closeErr)
		}
		db = nil
		return fmt.Errorf("error setting journal mode: %v", pragmaErr)
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

// addQuery adds the queries to the queries struct component
func (s *SDatabase) addQuery(slice *[]Query, query string, values []any) {
	*slice = append(*slice, Query{
		Query:  query,
		Values: values,
	})
}

// CreateTable creates a new table
func (s *SDatabase) CreateTable() {
	s.addQuery(&s.queries.CreateTable, sqlTables, nil)
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
	if joinClauses != nil {
		b.WriteString(" ")
		b.WriteString(buildJoinClauses(joinClauses))
	}

	// Build WHERE, GROUP BY, HAVING, ORDER BY, LIMIT, OFFSET
	b.WriteString(buildWhere(clauses.Where))
	if len(clauses.GroupBy) > 0 {
		// For SELECT, GroupBy is just: GROUP BY col1, col2 ...
		b.WriteString(fmt.Sprintf(" GROUP BY %s", strings.Join(clauses.GroupBy, ", ")))
	}
	if len(clauses.OrderBy) > 0 {
		b.WriteString(buildHaving(clauses.Having))
		b.WriteString(buildOrderBy(clauses.OrderBy))
	}
	if clauses.Limit != nil && *clauses.Limit > 0 {
		b.WriteString(buildLimit(int64(*clauses.Limit)))
	}
	if clauses.Offset != nil && *clauses.Offset > 0 {
		b.WriteString(buildOffset(int64(*clauses.Offset)))
	}
	b.WriteString(";")

	s.addQuery(&s.queries.Select, b.String(), nil)
}

// Delete deletes records from the table
func (s *SDatabase) Delete(table Table, clauses SelectClauses) {
	var b strings.Builder
	b.WriteString("DELETE FROM ")
	b.WriteString(table.Name)
	b.WriteString(buildWhere(clauses.Where))
	b.WriteString(";")

	s.addQuery(&s.queries.Delete, b.String(), nil)
}

// DeleteDb delete a db file
// - Used when the caching will be reset
func (s *SDatabase) DeleteDb() error {
	if ok, err := ValidateDbAbsPath(s.dbAbsPath); !ok {
		return err
	}

	// Proceed with deleting the database file
	err := osRemove(s.dbAbsPath)
	if err != nil {
		return err
	}

	return nil
}

// Apply executes all queued queries in s.queries
func (s *SDatabase) Apply() error {
	if !s.IsInitialized() {
		return fmt.Errorf(ErrorDatabaseNotInitialized)
	}

	// Begin a transaction
	sqlTx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Merge all queries in the desired order.
	var allQueries []Query
	allQueries = append(allQueries, s.queries.CreateTable...)
	allQueries = append(allQueries, s.queries.DropTable...)
	allQueries = append(allQueries, s.queries.Insert...)
	allQueries = append(allQueries, s.queries.Update...)
	allQueries = append(allQueries, s.queries.Delete...)
	allQueries = append(allQueries, s.queries.Select...)

	if len(allQueries) == 0 {
		return fmt.Errorf("error no queries")
	}

	err = s.exec(sqlTx, allQueries)
	if err != nil {
		return err
	}

	return nil
}

// QueryRows executes a SELECT query and returns results as a slice of maps.
func (s *SDatabase) QueryRows(query string, args ...any) ([]map[string]any, error) {
	if !s.IsInitialized() {
		return nil, fmt.Errorf(ErrorDatabaseNotInitialized)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %w", err)
	}

	var results []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if scanErr := rows.Scan(valuePtrs...); scanErr != nil {
			return nil, fmt.Errorf("error scanning row: %w", scanErr)
		}

		row := make(map[string]any, len(columns))
		for i, col := range columns {
			val := values[i]
			// Convert []byte to string for readability
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

// exec loops through all the queries and executes them
func (s *SDatabase) exec(sqlTx ISqlTx, allQueries []Query) error {
	for _, q := range allQueries {
		_, execErr := sqlTx.Exec(q.Query, q.Values...)
		if execErr != nil {
			rbErr := sqlTx.Rollback()
			if rbErr != nil {
				return fmt.Errorf("error rolling back transaction: %w (original error: %v)", rbErr, execErr)
			}
			return fmt.Errorf("error applying query (%s): %w", q.Query, execErr)
		}
	}

	if err := sqlTx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	s.queries = &Queries{
		CreateTable: []Query{},
		DropTable:   []Query{},
		Insert:      []Query{},
		Update:      []Query{},
		Delete:      []Query{},
		Select:      []Query{},
	}
	return nil
}
