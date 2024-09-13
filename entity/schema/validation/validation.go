package validation

import (
	"fmt"
	schemaPkg "github.com/AmadlaOrg/hery/entity/schema"
	"regexp"
	"strings"
)

// IValidation
type IValidation interface {
	Id(id, collectionName, entityUri string) error
}

// SValidation
type SValidation struct {
	Schema schemaPkg.ISchema
}

// Id validation of JSON-Schema for an entity
func (s *SValidation) Id(id, collectionName, entityUri string) error {
	if id == "" {
		return fmt.Errorf("schema validation failed: no `id` found in schema")
	}

	idUrnRegex := regexp.MustCompile(schemaPkg.EntityJsonSchemaIdURN)
	if !idUrnRegex.MatchString(id) {
		return fmt.Errorf("schema validation failed: invalid `id` found in schema")
	}

	prefix := s.Schema.GenerateURNPrefix(collectionName)
	if strings.HasPrefix(id, prefix) {
		return fmt.Errorf("schema validation failed: invalid `urn` found in schema")
	}

	if s.Schema.GenerateURN(prefix, entityUri) != id {
		return fmt.Errorf("schema validation failed: invalid `urn` found in schema")
	}

	return nil
}
