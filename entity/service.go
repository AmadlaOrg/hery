package entity

import "github.com/AmadlaOrg/hery/entity/version"

// NewEntityService to set up the entity build service
func NewEntityService() *SEntity {
	return &SEntity{
		EntityVersion: version.NewEntityVersionService(),

		// Data
		Entities: []Entity{},
	}
}
