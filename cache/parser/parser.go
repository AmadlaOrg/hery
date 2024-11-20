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

// Entity parses the entity to SQLite 3 struct that can be used in the query builder in the database package
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
		    id TEXT PRIMARY KEY NOT NULL DEFAULT
		);
	*/

	schema := entity.Schema.Schema

	// TODO:
	//entity.Schema.CompiledSchema.Types

	var (
		dynamicColumns       []database.Column
		dynamicRelationships []database.Relationships
	)
	for key, value := range schema {
		// Ensure "properties" is being processed
		if key == "properties" {
			properties, ok := value.(map[string]any)
			if !ok {
				continue // Skip if "properties" is not the expected type
			}

			for schemaPropertyName, schemaPropertyValue := range properties {
				// Assert schemaPropertyValue is a map
				propertyDetails, ok := schemaPropertyValue.(map[string]any)
				if !ok {
					continue // Skip if schemaPropertyValue is not a map
				}

				var dataType database.DataType

				// Change the data type if the "format" property is present
				if formatValue, ok := propertyDetails["format"].(string); ok {
					if dataFormat, valid := schemaStringToDataFormat(formatValue); valid {
						dataType = parseJsonSchemaFormatToSQLiteType(dataFormat)
					}
				} else if typeValue, ok := propertyDetails["type"].(string); ok {
					if dataTypeValue, valid := schemaStringToDataType(typeValue); valid {
						dataType = parseJsonSchemaToSQLiteType(dataTypeValue)
					}
				}

				// Append the column definition
				dynamicColumns = append(dynamicColumns, database.Column{
					ColumnName: schemaPropertyName,
					DataType:   dataType,
					Constraint: "", // TODO: Use constraints from JSON Schema (e.g., unique, required)
				})
			}
		}
	}

	entityBody := entity.Content.Body

	// TODO: It needs data type and constrain
	//var dynamicColumns []database.Column
	/*for key, value := range entityBody {
		//dataType := determineDataType(value)
		// TODO: Lookup the JsonSchema for the datatype

	}*/

	// TODO:
	/*var dynamicRelationships []database.Relationships
	for key, value := range entityBody {
		dynamicRelationships = append(dynamicRelationships, database.Relationships{})
	}*/

	return []database.Table{
		{},
		{},
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

//
// Private methods
//

func (s *SParser) databaseInsertTableEntities(entity entity.Entity) (*[]database.Table, error) {
	schema := entity.Schema.Schema
	// Convert schema map[string]any into a JSON string for cache storage
	schemaJsonBytes, err := jsonMarshal(schema)
	if err != nil {
		return nil, err
	}
	schemaJsonString := string(schemaJsonBytes)

	return &[]database.Table{
		{
			Name: "Entities",
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
			Name: "Ids",
			Rows: []map[string]any{
				{
					"Id":   entity.Id.String(),
					"Uri":  entity.Uri,
					"Name": entity.Name,
				},
			},
		},
	}, nil
}
