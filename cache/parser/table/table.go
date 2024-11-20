package table

import (
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/cache/parser/column"
)

var (
	Ids = database.Table{
		Name: "Ids",
		Columns: []database.Column{
			column.Id,
			{
				ColumnName: "CustomId",
				DataType:   "TEXT",
			},
			{
				ColumnName: "RefId",
				DataType:   "TEXT",
			},
			column.InsertDateTime,
			column.UpdateDateTime,
		},
	}
	Entities = database.Table{
		Name: "Entities",
		Columns: []database.Column{
			column.Id,
			{
				ColumnName: "Entity",
				DataType:   "TEXT",
			},
			{
				ColumnName: "Name",
				DataType:   "TEXT",
			},
			{
				ColumnName: "RepoUrl",
				DataType:   "TEXT",
			},
			{
				ColumnName: "Origin",
				DataType:   "TEXT",
			},
			{
				ColumnName: "Version",
				DataType:   "TEXT",
			},
			{
				ColumnName: "IsLatestVersion",
				DataType:   "BOOLEAN",
			},
			{
				ColumnName: "IsPseudoVersion",
				DataType:   "BOOLEAN",
			},
			{
				ColumnName: "AbsPath",
				DataType:   "TEXT",
			},
			{
				ColumnName: "Have",
				DataType:   "BOOLEAN",
			},
			{
				ColumnName: "Hash",
				DataType:   "TEXT",
			},
			{
				ColumnName: "Exist",
				DataType:   "BOOLEAN",
			},
			{
				ColumnName: "Schema",
				DataType:   "TEXT",
			},
			{
				ColumnName: "Content",
				DataType:   "TEXT",
			},
			column.InsertDateTime,
			column.UpdateDateTime,
		},
	}
)
