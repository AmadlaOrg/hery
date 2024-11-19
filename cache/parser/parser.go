package parser

import (
	"encoding/json"
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

var (
	jsonMarshal = json.Marshal
)

// Parse
func (s *SParser) Entity(entity entity.Entity) ([]database.Table, error) {
	// TODO: Use schema to determine the data type for the SQL
	// string == TEXT
	// TODO: Convert schema from the struct to the JSON-Schema string
	// TODO: For `Id` always: `Id TEXT PRIMARY KEY,`
	// TODO: Maybe always have `NOT NULL` as a constrain. E.g.: name TEXT NOT NULL

	// TODO: Handle different structures of _meta data
	// TODO: Single entity:
	/*
		_meta:
		  _entity: github.com/AmadlaOrg/Entity@latest
		  _body:
		    name: RandomName
		    description: Some description.
		    category: QA
	*/
	// TODO: Or list:
	/*
	  external-list:
	    - _entity: github.com/AmadlaOrg/QAFixturesSubEntityWithMultiSubEntities@latest
	      _body:
	        message: Another random message.
	    - _entity: github.com/AmadlaOrg/QAFixturesSubEntityWithMultiSubEntities@latest
	      _body:
	        message: Again, another random message.
	    - _entity: github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion@latest
	      _body:
	        title: Hello World!
	    - _entity: github.com/AmadlaOrg/QAFixturesEntityPseudoVersion@latest
	      _body:
	        name: John Doe
	*/

	// TODO: For UUID support
	/*
		CREATE TABLE example (
		    id TEXT PRIMARY KEY NOT NULL DEFAULT (lower(hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-4' || substr(hex(randomblob(2)), 2) || '-' || substr('89ab', abs(random() % 4) + 1, 1) || substr(hex(randomblob(2)), 2) || '-' || hex(randomblob(6))))
		);
	*/

	/*
		JSON-Schema Types:
		    string.
		    number.
		    integer.
		    object.
		    array.
		    boolean.
		    null.
	*/

	/*
		SQLite 3:
			Declared Type	Type Affinity	Storage Class
			INTEGER	INTEGER	INTEGER
			TINYINT, SMALLINT	INTEGER	INTEGER
			MEDIUMINT, BIGINT	INTEGER	INTEGER
			INT	INTEGER	INTEGER
			REAL, DOUBLE, FLOAT	REAL	REAL
			NUMERIC, DECIMAL	NUMERIC	REAL or INTEGER (if possible)
			TEXT	TEXT	TEXT
			CHARACTER, VARCHAR	TEXT	TEXT
			CLOB	TEXT	TEXT
			BLOB	BLOB	BLOB
			BOOLEAN	NUMERIC	INTEGER (1 for true, 0 for false)
			DATE, DATETIME	NUMERIC	TEXT, REAL, or INTEGER depending on the format
	*/

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

	// Convert schema map[string]any into a JSON string for cache storage
	schemaJsonBytes, err := jsonMarshal(entity.Schema.Schema)
	if err != nil {
		return nil, err
	}
	schemaJsonString := string(schemaJsonBytes)

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
					"Id":              entity.Id.String(),
					"Uri":             entity.Uri,
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
					"Schema":          schemaJsonString,
					"Content":         ``,
				},
			},
		},
		{
			Name: s.EntityToTableName(entity.Uri),
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
