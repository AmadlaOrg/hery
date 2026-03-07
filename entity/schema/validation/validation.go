package validation

import (
	"fmt"
	schemaPkg "github.com/AmadlaOrg/hery/entity/schema"
	"regexp"
	"strings"
)

// IValidation
type IValidation interface {
	Id(id, entityUri string) error
}

// SValidation
type SValidation struct {
	Schema schemaPkg.ISchema
}

// Id validation of JSON-Schema for an entity
func (s *SValidation) Id(id, entityUri string) error {

	// 1. Validates that the `id` is not empty
	if id == "" {
		return fmt.Errorf("schema validation failed: `id` is empty")
	}

	// 2. Validates that the format is good
	idUrnRegex := regexp.MustCompile(schemaPkg.EntityJsonSchemaIdURN)
	if !idUrnRegex.MatchString(id) {
		return fmt.Errorf("schema validation failed: invalid `id` format")
	}

	// 3. Validates that the prefix is standard (urn:hery:)
	if !strings.HasPrefix(id, "urn:hery:") {
		return fmt.Errorf("schema validation failed: invalid `urn` prefix (expected `urn:hery:`)")
	}

	// 4. Validates that the entire URN matches the entity URI
	expectedURN := s.Schema.GenerateURN(entityUri)
	if expectedURN != id {
		return fmt.Errorf(
			"schema validation failed: invalid `urn` — expected %s, got %s",
			expectedURN,
			id)
	}

	return nil
}
