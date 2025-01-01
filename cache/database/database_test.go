package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTable(t *testing.T) {
	//databaseService := NewDatabaseService()
	//databaseService.CreateTable()
}

func TestInsert(t *testing.T) {
	tests := []struct {
		name       string
		inputTable Table
	}{
		{
			name: "Test Insert",
			inputTable: Table{
				Name: "Net",
				Columns: []Column{
					{
						ColumnName: "Id",
						DataType:   "TEXT",
						Constraints: []Constraint{
							{
								Type: ConstraintPrimaryKey,
							},
						},
					},
					{
						ColumnName: "server_name",
						DataType:   "TEXT",
					},
					{
						ColumnName: "listen",
						DataType:   "TEXT",
					},
				},
				Rows: []map[string]any{
					{
						"Id":          "c6beaec1-90c4-4d2a-aaef-211ab00b86bd",
						"server_name": "localhost",
						"listen":      "[80, 443]",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			databaseService := NewDatabaseService()
			databaseService.Insert(tt.inputTable)
		})
	}
}

/*func TestInsert(t *testing.T) {
	// Arrange: Initialize the in-memory database
	dbPath := ":memory:"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize database service
	databaseService := &SDatabase{
		queries: &[]string{},
	}

	// Create table
	table := Table{
		Name: "Net",
		Columns: []Column{
			{
				ColumnName: "Id",
				DataType:   "TEXT",
				Constraints: []Constraint{
					{
						Type: ConstraintPrimaryKey,
					},
				},
			},
			{
				ColumnName: "server_name",
				DataType:   "TEXT",
			},
			{
				ColumnName: "listen",
				DataType:   "TEXT",
			},
		},
		Rows: []map[string]any{
			{
				"Id":          "c6beaec1-90c4-4d2a-aaef-211ab00b86bd",
				"server_name": "localhost",
				"listen":      "[80, 443]",
			},
		},
	}

	// Act: Create table and insert rows
	databaseService.CreateTable(table)
	for _, query := range *databaseService.queries {
		_, err := db.Exec(query)
		if err != nil {
			t.Fatalf("Failed to execute query: %v\nQuery: %s", err, query)
		}
	}

	databaseService.Insert(table)
	for _, query := range *databaseService.queries {
		_, err := db.Exec(query)
		if err != nil {
			t.Fatalf("Failed to execute query: %v\nQuery: %s", err, query)
		}
	}

	// Assert: Verify the data is inserted correctly
	var id, serverName, listen string
	err = db.QueryRow("SELECT Id, server_name, listen FROM Net WHERE Id = ?", "c6beaec1-90c4-4d2a-aaef-211ab00b86bd").Scan(&id, &serverName, &listen)
	if err != nil {
		t.Fatalf("Failed to query inserted row: %v", err)
	}

	if id != "c6beaec1-90c4-4d2a-aaef-211ab00b86bd" || serverName != "localhost" || listen != "[80, 443]" {
		t.Errorf("Inserted row does not match expected values: got (%s, %s, %s)", id, serverName, listen)
	}
}*/

func TestUpdate(t *testing.T) {
	tests := []struct {
		name       string
		inputTable Table
		expected   *Queries
	}{
		// TODO: It moves things around so sometimes the test pass other times no
		{
			name: "Test Insert",
			inputTable: Table{
				Name: "Net",
				Columns: []Column{
					{
						ColumnName: "Id",
						DataType:   "TEXT",
						Constraints: []Constraint{
							{
								Type: ConstraintPrimaryKey,
							},
						},
					},
					{
						ColumnName: "server_name",
						DataType:   "TEXT",
					},
					{
						ColumnName: "listen",
						DataType:   "TEXT",
					},
				},
				Rows: []map[string]any{
					{
						"Id":          "c6beaec1-90c4-4d2a-aaef-211ab00b86bd",
						"server_name": "localhost",
						"listen":      "[80, 443]",
					},
				},
			},
			expected: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update: []Query{
					{
						Query: "UPDATE Net SET Id = ?, server_name = ?, listen = ? WHERE Id = 'c6beaec1-90c4-4d2a-aaef-211ab00b86bd' AND server_name = 'localhost' AND listen = '[80, 443]'",
						Values: []string{
							"c6beaec1-90c4-4d2a-aaef-211ab00b86bd",
							"localhost",
							"[80, 443]",
						},
						Result: "",
					},
				},
				Delete: []Query{},
				Select: []Query{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			databaseService := SDatabase{
				queries: &Queries{
					CreateTable: []Query{},
					DropTable:   []Query{},
					Insert:      []Query{},
					Update:      []Query{},
					Delete:      []Query{},
					Select:      []Query{},
				},
			}
			databaseService.Update(tt.inputTable, map[string]any{
				"Id":          "c6beaec1-90c4-4d2a-aaef-211ab00b86bd",
				"server_name": "localhost",
				"listen":      "[80, 443]",
			})
			assert.Equal(t, tt.expected, databaseService.queries)
		})
	}
}

func TestDelete(t *testing.T) {}

func TestDropTable(t *testing.T) {}
