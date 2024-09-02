package cache

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// ICache
type ICache interface {
	Create()
	Insert()
	Select()
}

// SCache
type SCache struct{}

// TODO: journal mode: wal
// This will help with performance
// https://stackoverflow.com/questions/57118674/go-sqlite3-with-journal-mode-wal-gives-database-is-locked-error

// Create
func (s *SCache) Create() {
	// Open an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}(db)

	// Create a table
	createTableSQL := `CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	// Create an index on the name column
	createIndexSQL := `CREATE INDEX idx_name ON users(name);`
	_, err = db.Exec(createIndexSQL)
	if err != nil {
		fmt.Println("Error creating index:", err)
		return
	}
}

// Insert
func (s *SCache) Insert() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Unable to close database: %v", err)
		}
	}(db)

	// Insert some records
	insertUserSQL := `INSERT INTO users (name) VALUES (?)`
	for _, name := range []string{"Alice", "Bob", "Charlie"} {
		_, err = db.Exec(insertUserSQL, name)
		if err != nil {
			log.Fatalf("Error inserting record: %v", err)
			return
		}
	}
}

// Select
func (s *SCache) Select() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Unable to close database: %v", err)
		}
	}(db)

	// Query the record using the index
	var name string
	querySQL := `SELECT name FROM users WHERE name = 'Alice'`
	err = db.QueryRow(querySQL).Scan(&name)
	if err != nil {
		log.Fatalf("Error querying record: %v", err)
		return
	}

	// Print the result
	fmt.Println("User:", name)
}
