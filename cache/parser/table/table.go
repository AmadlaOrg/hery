package table

import (
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/cache/parser/column"
)

var (
	Entities = database.Table{
		Name: "Entities",
		Columns: []database.Column{
			column.Id,
			{
				ColumnName: "CustomId",
				DataType:   database.DataTypeText,
				Constraint: "",
			},
			{
				ColumnName: "Uri",
				DataType:   database.DataTypeText,
				Constraint: "",
			},
			{
				ColumnName: "Name",
				DataType:   database.DataTypeText,
				Constraint: "",
			},
			{
				ColumnName: "RepoUrl",
				DataType:   database.DataTypeText,
				Constraint: "",
			},
			{
				ColumnName: "Origin",
				DataType:   database.DataTypeText,
				Constraint: "",
			},
			{
				ColumnName: "Version",
				DataType:   database.DataTypeText,
				Constraint: "",
			},
			{
				ColumnName: "IsLatestVersion",
				DataType:   database.DataTypeBoolean,
				Constraint: "",
			},
			{
				ColumnName: "IsPseudoVersion",
				DataType:   database.DataTypeBoolean,
				Constraint: "",
			},
			{
				ColumnName: "AbsPath",
				DataType:   database.DataTypeText,
				Constraint: "",
			},
			{
				ColumnName: "Have",
				DataType:   database.DataTypeBoolean,
				Constraint: "",
			},
			{
				ColumnName: "Hash",
				DataType:   database.DataTypeText,
				Constraint: "",
			},
			{
				ColumnName: "Exist",
				DataType:   database.DataTypeBoolean,
				Constraint: "",
			},
			{
				ColumnName: "Schema",
				DataType:   database.DataTypeText,
				Constraint: "",
			},
			{
				ColumnName: "Content",
				DataType:   database.DataTypeText,
				Constraint: "",
			},
			column.InsertDateTime,
			column.UpdateDateTime,
		},
	}
)
