package collection

import "github.com/AmadlaOrg/hery/entity"

const (
	Match = `^[a-zA-Z0-9_-]+$`
)

// Collection contains all the collection data
type Collection struct {
	Name              string
	Paths             *[]string // TODO: storage.AbsStorage
	TransientEntities *[]*entity.Entity
}
