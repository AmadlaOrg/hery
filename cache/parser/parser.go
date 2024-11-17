package parser

import (
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/entity"
)

type IParser interface {
	ParseEntity(entity entity.Entity) (database.Table, error)
	ParseMultipleEntities(entities []entity.Entity) (database.Table, error)
	ParseTable(data []byte) (entity.Entity, error)
	ParseRow(data []byte) (entity.Entity, error)
}

type SParser struct{}

// Parse
func (s *SParser) ParseEntity(entity entity.Entity) (database.Table, error) {
	// TODO: Use schema to determine the data type for the SQL
	// string == TEXT
	// TODO: Convert schema from the struct to the JSON-Schema string
	// TODO: For `Id` always: `Id TEXT PRIMARY KEY,`
	// TODO: Maybe always have `NOT NULL` as a constrain. E.g.: name TEXT NOT NULL

	return database.Table{}, nil
}

// ParseMultipleEntities
func (s *SParser) ParseMultipleEntities(entities []entity.Entity) (database.Table, error) {
	return database.Table{}, nil
}

// ParseTable
func (s *SParser) ParseTable(data []byte) (entity.Entity, error) {
	return entity.Entity{}, nil
}

// ParseRow
func (s *SParser) ParseRow(data []byte) (entity.Entity, error) {
	return entity.Entity{}, nil
}
