package collection

import (
	"github.com/AmadlaOrg/hery/entity"
)

type IEntityCollection interface {
	setCollection(collectionName string) error
	Select() (entity.Entity, error)
	Add(entity entity.Entity) error
	Remove() error
}

type SEntityCollection struct {
}

// setCollection
func (s *SEntityCollection) setCollection(collectionName string) error {

	return nil
}

// Select entity
func (s *SEntityCollection) Select() (entity.Entity, error) {
	return entity.Entity{}, nil
}

// Add entity to a collection
func (s *SEntityCollection) Add(entity entity.Entity) error {

	return nil
}

// Remove entity
func (s *SEntityCollection) Remove() error {
	return nil
}
