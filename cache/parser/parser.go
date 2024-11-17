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
