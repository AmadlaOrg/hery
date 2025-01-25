package database

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/LibraryUtils/pointer"

	//"github.com/AmadlaOrg/hery/util/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Fixtures
var testDbAbsPath = "/tmp/hery.test.cache"

func TestInitialize_db_set(t *testing.T) {
	// Reset
	db = nil
	initialized = false

	/*mockSyncLocker := NewMockSyncLocker(t)
	mockSyncLocker.EXPECT().Lock()
	mockSyncLocker.EXPECT().Unlock()

	dbMutex = mockSyncLocker*/

	db = NewMockSqlDb(t)

	databaseService := NewDatabaseService(testDbAbsPath)
	err := databaseService.Initialize()

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

			/*mockSyncLocker := NewMockSyncLocker(t)
			mockSyncLocker.EXPECT().Lock()
			mockSyncLocker.EXPECT().Unlock()

			dbMutex = mockSyncLocker*/

			originalSqlOpen := sqlOpen
			defer func() { sqlOpen = originalSqlOpen }()
			sqlOpen = tt.internalSqlOpenFn

			databaseService := NewDatabaseService(testDbAbsPath)
			err := databaseService.Initialize()

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

			/*mockSyncLocker := NewMockSyncLocker(t)
			mockSyncLocker.EXPECT().Lock()
			mockSyncLocker.EXPECT().Unlock()

			dbMutex = mockSyncLocker*/

			databaseService := NewDatabaseService(testDbAbsPath)
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

			/*mockSyncLocker := NewMockSyncLocker(t)
			mockSyncLocker.EXPECT().Lock()
			mockSyncLocker.EXPECT().Unlock()

			dbMutex = mockSyncLocker*/

			databaseService := NewDatabaseService(testDbAbsPath)
			got := databaseService.IsInitialized()
			assert.Equal(t, tt.expectedInitialized, got)
		})
	}
}

func TestCreateTable(t *testing.T) {
	orginalSqlTables := sqlTables
	sqlTables = "SQL string to create table"
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
	databaseService.CreateTable()

	expected := &Queries{
		CreateTable: []Query{
			{
				Query: "SQL string to create table",
			},
		},
		DropTable: []Query{},
		Insert:    []Query{},
		Update:    []Query{},
		Delete:    []Query{},
		Select:    []Query{},
	}
	assert.Equal(t, databaseService.queries, expected)

	// Reset
	sqlTables = orginalSqlTables
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
			databaseService := NewDatabaseService(testDbAbsPath)
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
		/*{
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
		},*/
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

func TestSelect(t *testing.T) {
	tests := []struct {
		name             string
		inputTable       Table
		inputClauses     SelectClauses
		inputJoinClauses []JoinClauses
		expected         []Query
	}{
		{
			name: "Test Select: one condition",
			inputTable: Table{
				Name: "mock_table_name",
			},
			inputClauses: SelectClauses{
				Where: []Condition{
					{
						Column:   "server_name",
						Operator: OperatorLike,
						Value:    "%localhost",
					},
				},
			},
			expected: []Query{
				{
					Query: "SELECT * FROM mock_table_name WHERE server_name LIKE '%localhost';",
				},
			},
		},
		{
			name: "Test Select: two conditions",
			inputTable: Table{
				Name: "mock_table_name",
			},
			inputClauses: SelectClauses{
				Where: []Condition{
					{
						Column:   "server_name",
						Operator: OperatorLike,
						Value:    "%localhost",
					},
					{
						Column:   "foo",
						Operator: OperatorNotEqual,
						Value:    "boo",
					},
				},
			},
			expected: []Query{
				{
					Query: "SELECT * FROM mock_table_name WHERE server_name LIKE '%localhost' AND foo != 'boo';",
				},
			},
		},
		{
			name: "Test Select: two conditions and limit 10",
			inputTable: Table{
				Name: "mock_table_name",
			},
			inputClauses: SelectClauses{
				Where: []Condition{
					{
						Column:   "server_name",
						Operator: OperatorLike,
						Value:    "%localhost",
					},
					{
						Column:   "foo",
						Operator: OperatorNotEqual,
						Value:    "boo",
					},
				},
				Limit: pointer.ToPtr(10),
			},
			expected: []Query{
				{
					Query: "SELECT * FROM mock_table_name WHERE server_name LIKE '%localhost' AND foo != 'boo' LIMIT 10;",
				},
			},
		},
		{
			name: "Test Select: two conditions and limit 10 with offset 5",
			inputTable: Table{
				Name: "mock_table_name",
			},
			inputClauses: SelectClauses{
				Where: []Condition{
					{
						Column:   "server_name",
						Operator: OperatorLike,
						Value:    "%localhost",
					},
					{
						Column:   "foo",
						Operator: OperatorNotEqual,
						Value:    "boo",
					},
				},
				Limit:  pointer.ToPtr(10),
				Offset: pointer.ToPtr(5),
			},
			expected: []Query{
				{
					Query: "SELECT * FROM mock_table_name WHERE server_name LIKE '%localhost' AND foo != 'boo' LIMIT 10 OFFSET 5;",
				},
			},
		},
		{
			name: "Test Select: two conditions, and group by and limit 10 with offset 5",
			inputTable: Table{
				Name: "mock_table_name",
			},
			inputClauses: SelectClauses{
				Where: []Condition{
					{
						Column:   "server_name",
						Operator: OperatorLike,
						Value:    "%localhost",
					},
					{
						Column:   "foo",
						Operator: OperatorNotEqual,
						Value:    "boo",
					},
				},
				GroupBy: []string{
					"server_name",
				},
				Limit:  pointer.ToPtr(10),
				Offset: pointer.ToPtr(5),
			},
			expected: []Query{
				{
					Query: "SELECT * FROM mock_table_name WHERE server_name LIKE '%localhost' AND foo != 'boo' GROUP BY server_name LIMIT 10 OFFSET 5;",
				},
			},
		},
		{
			name: "Test Select: two conditions, and group by and limit 10 with offset 5",
			inputTable: Table{
				Name: "mock_table_name",
			},
			inputClauses: SelectClauses{
				Where: []Condition{
					{
						Column:   "server_name",
						Operator: OperatorLike,
						Value:    "%localhost",
					},
					{
						Column:   "foo",
						Operator: OperatorNotEqual,
						Value:    "boo",
					},
				},
				GroupBy: []string{
					"server_name",
				},
				Limit:  pointer.ToPtr(10),
				Offset: pointer.ToPtr(5),
			},
			expected: []Query{
				{
					Query: "SELECT * FROM mock_table_name WHERE server_name LIKE '%localhost' AND foo != 'boo' GROUP BY server_name LIMIT 10 OFFSET 5;",
				},
			},
		},
		{
			name: "Test Select: two conditions, and group by, order by and limit 10 with offset 5",
			inputTable: Table{
				Name: "mock_table_name",
			},
			inputClauses: SelectClauses{
				Where: []Condition{
					{
						Column:   "server_name",
						Operator: OperatorLike,
						Value:    "%localhost",
					},
					{
						Column:   "foo",
						Operator: OperatorNotEqual,
						Value:    "boo",
					},
				},
				GroupBy: []string{
					"server_name",
				},
				OrderBy: []OrderBy{
					{
						Column:    "server_name",
						Direction: OrderByDesc,
					},
				},
				Limit:  pointer.ToPtr(10),
				Offset: pointer.ToPtr(5),
			},
			expected: []Query{
				{
					Query: "SELECT * FROM mock_table_name WHERE server_name LIKE '%localhost' AND foo != 'boo' GROUP BY server_name ORDER BY server_name DESC LIMIT 10 OFFSET 5;",
				},
			},
		},
		/*{
			name: "Test Select: inner joins",
			inputTable: Table{
				Name: "mock_table_name",
			},
			inputClauses: SelectClauses{},
			inputJoinClauses: []JoinClauses{
				{
					Table: "mock_second_table_name",
					Type:  JoinTypeInner,
					On: []JoinCondition{
						{
							Column1:  "server_name",
							Operator: OperatorLike,
							Column2:  "local%",
						},
					},
				},
			},
			expected: []Query{
				{
					Query: "SELECT * FROM mock_table_name INNER JOIN mock_second_table_name ON mock_second_table_name.server_name LIKE mock_table_name.server_name;",
				},
			},
		},*/
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
			databaseService.Select(tt.inputTable, tt.inputClauses, tt.inputJoinClauses)
			assert.Equal(t, tt.expected, databaseService.queries.Select)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name         string
		inputTable   Table
		inputClauses SelectClauses
		expected     *Queries
	}{
		{
			name: "Simple WHERE clause with column selection",
			inputTable: Table{
				Name: "mock_table_name",
			},
			inputClauses: SelectClauses{
				Where: []Condition{
					{
						Column:   "Id",
						Operator: OperatorEqual,
						Value:    "mock_table_id",
					},
				},
			},
			expected: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update:      []Query{},
				Delete: []Query{
					{
						Query: "DELETE FROM mock_table_name WHERE Id = 'mock_table_id';",
					},
				},
				Select: []Query{},
			},
		},
		{
			name: "Simple WHERE clause with TWO column selection",
			inputTable: Table{
				Name: "mock_table_name",
			},
			inputClauses: SelectClauses{
				Where: []Condition{
					{
						Column:   "Name",
						Operator: OperatorEqual,
						Value:    "mock_table_column_name",
					},
					{
						Column:   "Type",
						Operator: OperatorLike,
						Value:    "%mock_table_column_type%",
					},
				},
			},
			expected: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update:      []Query{},
				Delete: []Query{
					{
						Query: "DELETE FROM mock_table_name WHERE Name = 'mock_table_column_name' AND Type LIKE '%mock_table_column_type%';",
					},
				},
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
			databaseService.Delete(tt.inputTable, tt.inputClauses)
			assert.Equal(t, tt.expected, databaseService.queries)
		})
	}
}

func TestDeleteDb(t *testing.T) {
	tests := []struct {
		name                     string
		serviceStructDbAbsPath   string
		internalFileIsValidMagic func(path string, magic []byte) (bool, error)
		internalOsRemove         func(name string) error
		expectedError            error
		hasError                 bool
	}{
		{
			name:                   "Success Delete",
			serviceStructDbAbsPath: "testdata/mock_db",
			internalFileIsValidMagic: func(path string, magic []byte) (bool, error) {
				return true, nil
			},
			internalOsRemove: func(name string) error {
				return nil
			},
			hasError: false,
		},
		//
		// Errors
		//
		{
			name:                   "Error: Delete fails at ValidateDbAbsPath",
			serviceStructDbAbsPath: "testdata/mock_db",
			internalFileIsValidMagic: func(path string, magic []byte) (bool, error) {
				return false, errors.New("test error (ValidateDbAbsPath)")
			},
			expectedError: errors.New("test error (ValidateDbAbsPath)"),
			hasError:      true,
		},
		{
			name:                   "Error: Delete fails at os.Remove",
			serviceStructDbAbsPath: "testdata/mock_db",
			internalFileIsValidMagic: func(path string, magic []byte) (bool, error) {
				return true, nil
			},
			internalOsRemove: func(name string) error {
				return errors.New("test error (OsRemove)")
			},
			expectedError: errors.New("test error (OsRemove)"),
			hasError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			databaseService := SDatabase{
				dbAbsPath: tt.serviceStructDbAbsPath,
				queries: &Queries{
					CreateTable: []Query{},
					DropTable:   []Query{},
					Insert:      []Query{},
					Update:      []Query{},
					Delete:      []Query{},
					Select:      []Query{},
				},
			}

			originalFileIsValidMagic := fileIsValidMagic
			defer func() { fileIsValidMagic = originalFileIsValidMagic }()
			fileIsValidMagic = tt.internalFileIsValidMagic

			originalOsRemove := osRemove
			defer func() { osRemove = originalOsRemove }()
			osRemove = tt.internalOsRemove

			err := databaseService.DeleteDb()
			if tt.hasError {
				assert.Error(t, err)
				assert.ErrorContains(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApply(t *testing.T) {
	tests := []struct {
		name                     string
		internalQueries          *Queries
		internalInitialized      bool
		internalDbBeginExpect    bool
		internalDbBeginError     error
		internalTxExecExpect     bool
		internalTxExecError      error
		internalTxRollbackExpect bool
		internalTxRollbackError  error
		internalTxCommitExpect   bool
		internalTxCommitError    error
		expectedQueries          *Queries
		expectedError            error
		hasError                 bool
	}{
		{
			name: "Error: is not initialized",
			internalQueries: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update:      []Query{},
				Delete:      []Query{},
				Select:      []Query{},
			},
			internalInitialized: false,
			expectedQueries: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update:      []Query{},
				Delete:      []Query{},
				Select:      []Query{},
			},
			expectedError: fmt.Errorf(ErrorDatabaseNotInitialized),
			hasError:      true,
		},
		{
			name: "Error: db Begin fails",
			internalQueries: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update:      []Query{},
				Delete:      []Query{},
				Select:      []Query{},
			},
			internalInitialized:   true,
			internalDbBeginExpect: true,
			internalDbBeginError:  errors.New("test error (ValidateDbBegin)"),
			expectedQueries: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update:      []Query{},
				Delete:      []Query{},
				Select:      []Query{},
			},
			expectedError: fmt.Errorf("test error (ValidateDbBegin)"),
			hasError:      true,
		},
		{
			name: "Error: no queries",
			internalQueries: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update:      []Query{},
				Delete:      []Query{},
				Select:      []Query{},
			},
			internalInitialized:   true,
			internalDbBeginExpect: true,
			internalDbBeginError:  nil,
			expectedQueries: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update:      []Query{},
				Delete:      []Query{},
				Select:      []Query{},
			},
			expectedError: fmt.Errorf("error no queries"),
			hasError:      true,
		},
		// TODO: Issue with interface and struct
		/*{
			name: "Error: tx Exec fails",
			internalQueries: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update:      []Query{},
				Delete:      []Query{},
				Select: []Query{
					{
						Query: "select * from table1",
					},
				},
			},
			internalInitialized:   true,
			internalDbBeginExpect: true,
			internalDbBeginError:  nil,
			internalTxExecExpect:  true,
			internalTxExecError:   errors.New("test error (tx Exec)"),
			expectedQueries: &Queries{
				CreateTable: []Query{},
				DropTable:   []Query{},
				Insert:      []Query{},
				Update:      []Query{},
				Delete:      []Query{},
				Select: []Query{
					{
						Query: "select * from table1",
					},
				},
			},
			expectedError: fmt.Errorf("test error (tx Exec)"),
			hasError:      true,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			databaseService := &SDatabase{
				queries: tt.internalQueries,
			}

			/*mockSyncLocker := NewMockSyncLocker(t)
			mockSyncLocker.EXPECT().Lock()
			mockSyncLocker.EXPECT().Unlock()

			dbMutex = mockSyncLocker*/

			// initialized
			originalInitialized := initialized
			defer func() { initialized = originalInitialized }()
			initialized = tt.internalInitialized

			mockSqlDb := NewMockSqlDb(t)

			if tt.internalDbBeginExpect {

				if tt.internalTxExecExpect {
					mockSqlTx := NewMockSqlTx(t)
					mockSqlTx.EXPECT().Exec(mock.Anything).Return(nil, tt.internalTxExecError)

					if tt.internalTxRollbackExpect {
						mockSqlTx.EXPECT().Rollback().Return(tt.internalTxRollbackError)
					}

					if tt.internalTxCommitExpect {
						mockSqlTx.EXPECT().Commit().Return(tt.internalTxCommitError)
					}

					// FIXME:
					//mockSqlDb.EXPECT().Begin().Return(mockSqlTx, tt.internalDbBeginError)
				} else {
					mockSqlDb.EXPECT().Begin().Return(nil, tt.internalDbBeginError)
				}

				/*originalTxExec := txExec
				defer func() { txExec = originalTxExec }()
				txExec = func(tx *sql.Tx, query string, args ...any) (sql.Result, error) {
					return mockSqlTx.Exec(query, args...)
				}*/

				db = mockSqlDb
			}

			err := databaseService.Apply()
			if tt.hasError {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedQueries, databaseService.queries)
		})
	}
}
