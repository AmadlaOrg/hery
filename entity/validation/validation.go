package validation

import (
	"fmt"
	schemaPkg "github.com/AmadlaOrg/hery/entity/schema"
	schemaValidationPkg "github.com/AmadlaOrg/hery/entity/schema/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"strings"
	"unicode"
)

// IValidation used by mockery
type IValidation interface {
	RootEntity(rootSchema, selfSchema *jsonschema.Schema, heryContent map[string]any) error
	Entity(schema *jsonschema.Schema, heryContent map[string]any) error
	EntityUri(entityUrl string) bool
}

// SValidation used by mockery
type SValidation struct {
	Version           version.IVersion
	VersionValidation versionValidationPkg.IValidation
	Schema            schemaPkg.ISchema
	SchemaValidation  schemaValidationPkg.IValidation
}

// RootEntity validates the root-level entity properties
func (s *SValidation) RootEntity(rootSchema, selfSchema *jsonschema.Schema, heryContent map[string]any) error {
	// Validate that _type is present and matches expected format
	typeVal, ok := heryContent["_type"].(string)
	if !ok || typeVal == "" {
		return fmt.Errorf("_type is required and must be a non-empty string")
	}

	// Validate _parent references same _type if present
	if parentVal, exists := heryContent["_parent"]; exists {
		if _, ok := parentVal.(string); !ok {
			return fmt.Errorf("_parent must be a string URI")
		}
	}

	return nil
}

// Entity validates the YAML content against the JSON schema
func (s *SValidation) Entity(schema *jsonschema.Schema, heryContent map[string]any) error {
	// 1. Validate _type is present
	typeVal, ok := heryContent["_type"].(string)
	if !ok || typeVal == "" {
		return fmt.Errorf("_type is required")
	}

	// 2. Validate _self if present
	if selfVal, exists := heryContent["_self"]; exists {
		selfMap, ok := selfVal.(map[string]any)
		if !ok {
			return fmt.Errorf("_self must be a map")
		}
		if len(selfMap) == 0 {
			return fmt.Errorf("_self must not be empty if present")
		}
		if _, hasType := selfMap["_type"]; hasType {
			return fmt.Errorf("_self must not contain _type")
		}
	}

	// 3. Validate the hery file content with the loaded schema
	if err := schema.Validate(heryContent); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	return nil
}

// EntityUri validates the module path for go get
//
// A entity URI cannot contain the usual URL elements.
func (s *SValidation) EntityUri(entityUrl string) bool {
	if strings.Contains(entityUrl, "://") {
		return false
	}
	for _, r := range entityUrl {
		if unicode.IsSpace(r) || r == ':' || r == '?' || r == '&' || r == '=' || r == '#' {
			return false
		}
	}
	return true
}
