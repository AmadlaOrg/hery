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
// Multiple layers of validation helps the developer debug with more specific errors
func (s *SValidation) Id(id, collectionName, entityUri string) error {

	// 1. Validates that the `id` is not empty
	if id == "" {
		return fmt.Errorf("schema validation failed: `id` is empty or not set")
	}

	// 2. Validates that the format is good
	idUrnRegex := regexp.MustCompile(schemaPkg.EntityJsonSchemaIdURN)
	if !idUrnRegex.MatchString(id) {
		return fmt.Errorf("schema validation failed: invalid `id` format")
	}

	// 3. Validates that the prefix of the HERY URN is standard
	prefix := s.Schema.GenerateURNPrefix(collectionName)
	if strings.HasPrefix(id, prefix) {
		return fmt.Errorf("schema validation failed: invalid `urn` prefix (`urn:hery:<collection name>:`)")
	}

	// 4. Validates that the entire URN is valid
	if s.Schema.GenerateURN(prefix, entityUri) != id {
		return fmt.Errorf("schema validation failed: invalid `urn` since the other values are not the same as the `_entity` URI")
	}

	return nil
}
