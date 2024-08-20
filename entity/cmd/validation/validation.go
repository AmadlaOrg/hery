package validation

import (
	"errors"
)

// Entities function is for validating passing entities in args
func Entities(entities []string) error {
	numEntities := len(entities)
	if numEntities == 0 {
		return errors.New("no entity URI specified")
	} else if numEntities > 60 {
		return errors.New("too many entity URIs (the limit is 60)")
	}
	return nil
}
