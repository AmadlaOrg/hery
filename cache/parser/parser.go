package parser

import (
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/entity"
	"regexp"
	"strings"
)

type IParser interface {
	Entity(entity entity.Entity) ([]database.Table, error)
	EntityToTableName(entity string) string
	DatabaseTable(data []byte) (entity.Entity, error)
	DatabaseRow(data []byte) (entity.Entity, error)
}

type SParser struct{}

// Parse
func (s *SParser) Entity(entity entity.Entity) ([]database.Table, error) {
	// TODO: Use schema to determine the data type for the SQL
	// string == TEXT
	// TODO: Convert schema from the struct to the JSON-Schema string
	// TODO: For `Id` always: `Id TEXT PRIMARY KEY,`
	// TODO: Maybe always have `NOT NULL` as a constrain. E.g.: name TEXT NOT NULL

	entityBody := entity.Content.Body

	// TODO: It needs data type and constrain
	var dynamicColumns []database.Column
	for key, value := range entityBody {
		//dataType := determineDataType(value)
		// TODO: Lookup the JsonSchema for the datatype
		dynamicColumns = append(dynamicColumns, database.Column{
			ColumnName: key,
			DataType:   "", //dataType,
			Constraint: "", // TODO: Maybe there is something we can use from json schema (unique, require, etc)
		})
	}

	// TODO:
	var dynamicRelationships []database.Relationships
	for key, value := range entityBody {
		dynamicRelationships = append(dynamicRelationships, database.Relationships{})
	}

	return []database.Table{
		{
			Name: "Entities",
			Columns: []database.Column{
				{
					ColumnName: "Id",
					DataType:   "string",
				},
				{
					ColumnName: "Entity",
					DataType:   "string",
				},
				{
					ColumnName: "Name",
					DataType:   "string",
				},
				{
					ColumnName: "RepoUrl",
					DataType:   "string",
				},
				{
					ColumnName: "Origin",
					DataType:   "string",
				},
				{
					ColumnName: "Version",
					DataType:   "string",
				},
				{
					ColumnName: "IsLatestVersion",
					DataType:   "bool",
				},
				{
					ColumnName: "IsPseudoVersion",
					DataType:   "bool",
				},
				{
					ColumnName: "AbsPath",
					DataType:   "string",
				},
				{
					ColumnName: "Have",
					DataType:   "bool",
				},
				{
					ColumnName: "Hash",
					DataType:   "string",
				},
				{
					ColumnName: "Exist",
					DataType:   "bool",
				},
				{
					ColumnName: "Schema",
					DataType:   "string",
				},
				{
					ColumnName: "Content",
					DataType:   "string",
				},
			},
			Rows: []map[string]any{
				{
					"Id":              entity.Id,
					"Entity":          entity.Entity,
					"Name":            entity.Name,
					"RepoUrl":         entity.RepoUrl,
					"Origin":          entity.Origin,
					"Version":         entity.Version,
					"IsLatestVersion": entity.IsLatestVersion,
					"IsPseudoVersion": entity.IsPseudoVersion,
					"AbsPath":         entity.AbsPath,
					"Have":            entity.Have,
					"Hash":            entity.Hash,
					"Exist":           entity.Exist,
					"Schema":          ``,
					"Content":         ``,
				},
			},
		},
		{
			Name: s.EntityToTableName(entity.Entity),
			Columns: []database.Column{
				{},
			},
		},
	}, nil
}

// EntityToTableName
func (s *SParser) EntityToTableName(entity string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	tableName := re.ReplaceAllString(entity, "_")
	return strings.Trim(tableName, "_")
}

// ParseTable
func (s *SParser) DatabaseTable(data []byte) (entity.Entity, error) {
	return entity.Entity{}, nil
}

// ParseRow
func (s *SParser) DatabaseRow(data []byte) (entity.Entity, error) {
	return entity.Entity{}, nil
}
