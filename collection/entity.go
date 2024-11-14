package collection

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
)

type IEntityCollection interface {
}

type SEntityCollection struct {
}

func (s *SEntityCollection) setCollection(collectionName string) error {

	return nil
}

// Select
func (s *SEntityCollection) Select(collectionName string) (entity.Entity, error) {
	return entity.Entity{}, nil
}

// Add to a collection
func (s *SEntityCollection) Add(collectionName string, entity entity.Entity) error {
	if collectionName == "" {
		return fmt.Errorf("collection name is empty")
	}

	return nil
}

// Remove
func (s *SEntityCollection) Remove(collectionName string) error {
	return nil
}
