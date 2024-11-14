package collection

import "github.com/AmadlaOrg/hery/entity"

const (
	Match = `^[a-zA-Z0-9_-]+$`
)

// Collection contains all the collection data
type Collection struct {
	Name     string
	Entities []entity.Entity
	Paths    []string // TODO: storage.AbsStorage
}

type EntityCollection struct {
}
