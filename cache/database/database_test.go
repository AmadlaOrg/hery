package database

import (
	"testing"
)

func TestCreateTable(t *testing.T) {
	//databaseService := NewDatabaseService()
	//databaseService.CreateTable()
}

func TestInsert(t *testing.T) {
	databaseService := NewDatabaseService()
	databaseService.Insert(Table{
		Name: "Net",
		Columns: []Column{
			{
				ColumnName: "Id",
				DataType:   "string",
			},
			{
				ColumnName: "server_name",
				DataType:   "string",
			},
			{
				ColumnName: "listen",
				DataType:   "string",
			},
		},
		Rows: []map[string]any{
			{
				"Id":    "c6beaec1-90c4-4d2a-aaef-211ab00b86bd",
				"ports": "[80, 443]",
			},
		},
	})
}

func TestUpdate(t *testing.T) {}

func TestDelete(t *testing.T) {}

func TestDropTable(t *testing.T) {}
