package validation

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
	schemaPkg "github.com/AmadlaOrg/hery/entity/schema"
	"regexp"
	"strings"
)

// IValidation
type IValidation interface {
	Id(id, collectionName string, entityMeta entity.Entity) error
}

// SValidation
type SValidation struct{}

// Id validation of JSON-Schema for an entity
func (s *SValidation) Id(id, collectionName string, entityMeta entity.Entity) error {
	if id == "" {
		return fmt.Errorf("schema validation failed: no `id` found in schema")
	}

	idUrnRegex := regexp.MustCompile(schemaPkg.EntityJsonSchemaIdURN)
	if !idUrnRegex.MatchString(id) {
		return fmt.Errorf("schema validation failed: invalid `id` found in schema")
	}

	prefix := fmt.Sprintf("urn:hery:%s:", collectionName)
	if strings.HasPrefix(id, prefix) {
		return fmt.Errorf("schema validation failed: invalid `urn` found in schema")
	}

	suffix := fmt.Sprintf(":%s:%s", entityMeta.Name, entityMeta.Version)
	if strings.HasSuffix(id, suffix) {
		return fmt.Errorf("schema validation failed: invalid `urn` found in schema")
	}

	return nil
}
