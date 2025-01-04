package database

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestInitialize_db_set(t *testing.T) {
	// Reset
	db = nil
	initialized = false

	mockSyncLocker := NewMockSyncLocker(t)
	mockSyncLocker.EXPECT().Lock()
	mockSyncLocker.EXPECT().Unlock()

	dbMutex = mockSyncLocker
	db = NewMockSqlDb(t)

	databaseService := NewDatabaseService()
	err := databaseService.Initialize("/tmp/hery.test.cache")

	assert.NoError(t, err)
}

func TestInitialize(t *testing.T) {
	tests := []struct {
		name              string
		inputDbPath       string
		internalSqlOpenFn func(driverName, dataSourceName string) (ISqlDb, error)
		expectedErr       error
		hasError          bool
	}{
		{
			name:        "Initialize database",
			inputDbPath: "/tmp/hery.test.cache",
			internalSqlOpenFn: func(driverName, dataSourceName string) (ISqlDb, error) {
				mockSqlDb := NewMockSqlDb(t)
				mockSqlDb.EXPECT().Exec(mock.Anything).Return(nil, nil)
				mockSqlDb.EXPECT().SetMaxOpenConns(mock.Anything)
				mockSqlDb.EXPECT().SetMaxIdleConns(mock.Anything)
				mockSqlDb.EXPECT().SetConnMaxLifetime(mock.Anything)
				return mockSqlDb, nil
			},
			hasError: false,
		},
		//
		// Error
		//
		{
			name:        "Error: db Open fail",
			inputDbPath: "/tmp/hery.test.cache",
			internalSqlOpenFn: func(driverName, dataSourceName string) (ISqlDb, error) {
				return nil, assert.AnError
			},
			expectedErr: errors.New("error opening database: "),
			hasError:    true,
		},
		{
			name:        "Error: db Exec function throws error",
			inputDbPath: "/tmp/hery.test.cache",
			internalSqlOpenFn: func(driverName, dataSourceName string) (ISqlDb, error) {
				mockSqlDb := NewMockSqlDb(t)
				mockSqlDb.EXPECT().Exec(mock.Anything).Return(nil, assert.AnError) // assert.AnError
				mockSqlDb.EXPECT().Close().Return(nil)
				return mockSqlDb, nil
			},
			expectedErr: errors.New("error setting journal mode: "),
			hasError:    true,
		},
		{
			name:        "Error: db Exec function throws error also db Close",
			inputDbPath: "/tmp/hery.test.cache",
			internalSqlOpenFn: func(driverName, dataSourceName string) (ISqlDb, error) {
				mockSqlDb := NewMockSqlDb(t)
				mockSqlDb.EXPECT().Exec(mock.Anything).Return(nil, assert.AnError) // assert.AnError
				mockSqlDb.EXPECT().Close().Return(assert.AnError)
				return mockSqlDb, nil
			},
			expectedErr: assert.AnError,
			hasError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset
			db = nil
			initialized = false

			mockSyncLocker := NewMockSyncLocker(t)
			mockSyncLocker.EXPECT().Lock()
			mockSyncLocker.EXPECT().Unlock()

			dbMutex = mockSyncLocker

			originalSqlOpen := sqlOpen
			defer func() { sqlOpen = originalSqlOpen }()
			sqlOpen = tt.internalSqlOpenFn

			databaseService := NewDatabaseService()
			err := databaseService.Initialize("/tmp/hery.test.cache")

			if tt.hasError {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClose(t *testing.T) {
	tests := []struct {
		name                string
		externalInitialized bool
		externalDb          ISqlDb
		hasError            bool
	}{
		{
			name:                "Close database",
			externalInitialized: true,
			externalDb: func() ISqlDb {
				mockSqlDb := NewMockSqlDb(t)
				mockSqlDb.EXPECT().Close().Return(nil)
				return mockSqlDb
			}(),
			hasError: false,
		},
		{
			name:                "Database is closed",
			externalInitialized: true,
			externalDb: func() ISqlDb {
				return nil
			}(),
			hasError: false,
		},
		//
		// Error
		//
		{
			name:                "Error: Close database",
			externalInitialized: true,
			externalDb: func() ISqlDb {
				mockSqlDb := NewMockSqlDb(t)
				mockSqlDb.EXPECT().Close().Return(assert.AnError)
				return mockSqlDb
			}(),
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db = tt.externalDb
			initialized = tt.externalInitialized

			mockSyncLocker := NewMockSyncLocker(t)
			mockSyncLocker.EXPECT().Lock()
			mockSyncLocker.EXPECT().Unlock()

			dbMutex = mockSyncLocker

			databaseService := NewDatabaseService()
			err := databaseService.Close()

			if tt.hasError {
				assert.Error(t, err)
				assert.True(t, initialized)
				assert.NotNil(t, db)
			} else {
				assert.NoError(t, err)
				assert.False(t, initialized)
				assert.Nil(t, db)
			}

		})
	}
}

func TestIsInitialized(t *testing.T) {
	tests := []struct {
		name                string
		externalInitialized bool
		expectedInitialized bool
	}{
		{
			name:                "IsInitialized is true",
			externalInitialized: true,
			expectedInitialized: true,
		},
		{
			name:                "IsInitialized is false",
			externalInitialized: false,
			expectedInitialized: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialized = tt.externalInitialized

			mockSyncLocker := NewMockSyncLocker(t)
			mockSyncLocker.EXPECT().Lock()
			mockSyncLocker.EXPECT().Unlock()

			dbMutex = mockSyncLocker

			databaseService := NewDatabaseService()
			got := databaseService.IsInitialized()
			assert.Equal(t, tt.expectedInitialized, got)
		})
	}
}

func TestCreateTable(t *testing.T) {
	tests := []struct {
		name       string
		inputTable Table
		expected   *Queries
	}{
		{
			name:       "Create table successfully",
			inputTable: Table{},
			expected:   &Queries{},
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
			databaseService.CreateTable(tt.inputTable)
			assert.Equal(t, tt.expected, databaseService.queries)
		})
	}
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
			databaseService.Update(tt.inputTable, []Condition{
				{Column: "Id", Operator: "=", Value: "c6beaec1-90c4-4d2a-aaef-211ab00b86bd"},
				{Column: "server_name", Operator: "LIKE", Value: "localhost"},
				{Column: "listen", Operator: "IN", Value: "[80, 443]"},
			})
			assert.Equal(t, tt.expected, databaseService.queries)
		})
	}
}

func TestDelete(t *testing.T) {}

func TestDropTable(t *testing.T) {}
