package cache

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Cache interface {
	Create()
	Insert()
	Select()
}

func Create() {
	// Open an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

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

func Insert() {
	db, err := sql.Open("sqlite3", ":memory:")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Insert some records
	insertUserSQL := `INSERT INTO users (name) VALUES (?)`
	for _, name := range []string{"Alice", "Bob", "Charlie"} {
		_, err = db.Exec(insertUserSQL, name)
		if err != nil {
			fmt.Println("Error inserting record:", err)
			return
		}
	}
}

func Select() {
	db, err := sql.Open("sqlite3", ":memory:")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Query the record using the index
	var name string
	querySQL := `SELECT name FROM users WHERE name = 'Alice'`
	err = db.QueryRow(querySQL).Scan(&name)
	if err != nil {
		fmt.Println("Error querying record:", err)
		return
	}

	// Print the result
	fmt.Println("User:", name)
}
