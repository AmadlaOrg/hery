package database

import (
	"database/sql"
	"fmt"
	"log"
)

// IDatabase
type IDatabase interface {
	CreateTable(table Table) error
	Insert(table Table) error
	Update(table Table) error
	Select(table Table) (string, error)
	Delete(table Table) error
	DropTable(table Table) error
}

// SDatabase
type SDatabase struct{}

var (
	db *sql.DB
)

func init() {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(fmt.Errorf("error opening database: %v", err))
	}
	db = conn
}

// TODO: journal mode: wal
// This will help with performance
// https://stackoverflow.com/questions/57118674/go-sqlite3-with-journal-mode-wal-gives-database-is-locked-error

// Create
func (s *SDatabase) CreateTable(table Table) error {
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}(db)

	// Create a table
	createTableSQL := fmt.Sprintf(`CREATE TABLE %s (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    );`, table.Name)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	// Create an index on the name column
	createIndexSQL := `CREATE INDEX idx_name ON users(name);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		return fmt.Errorf("error creating index: %v", err)
	}

	return nil
}

// Insert
func (s *SDatabase) Insert(table Table) error {
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Unable to close database: %v", err)
		}
	}(db)

	// Insert some records
	insertUserSQL := `INSERT INTO users (name) VALUES (?)`
	for _, name := range []string{"Alice", "Bob", "Charlie"} {
		_, err := db.Exec(insertUserSQL, name)
		if err != nil {
			return fmt.Errorf("error inserting record: %v", err)
		}
	}

	return nil
}

func (s *SDatabase) Update(table Table) error {
	return nil
}

// Select
func (s *SDatabase) Select(table Table) (string, error) {
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Unable to close database: %v", err)
		}
	}(db)

	// Query the record using the index
	var name string
	querySQL := `SELECT name FROM users WHERE name = 'Alice'`
	err := db.QueryRow(querySQL).Scan(&name)
	if err != nil {
		return "", fmt.Errorf("error querying record: %v", err)
	}

	return name, nil
}

func (s *SDatabase) Delete(table Table) error {

	return nil
}

func (s *SDatabase) DropTable(table Table) error {

	return nil
}
