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

// Validator used by mockery
type Validator interface {
	RootEntity(rootSchema, selfSchema *jsonschema.Schema, heryContent map[string]any) error
	Entity(schema *jsonschema.Schema, heryContent map[string]any) error
	EntityUri(entityUrl string) bool
}

// validator used by mockery
type validator struct {
	Version           version.Version
	VersionValidation versionValidationPkg.Validator
	Schema            schemaPkg.Schema
	SchemaValidation  schemaValidationPkg.Validator
}

// RootEntity validates the root-level entity properties
func (s *validator) RootEntity(rootSchema, selfSchema *jsonschema.Schema, heryContent map[string]any) error {
	// Validate that _type is present and matches expected format
	typeVal, ok := heryContent["_type"].(string)
	if !ok || typeVal == "" {
		return fmt.Errorf("_type is required and must be a non-empty string")
	}

	// Validate _extends references same _type if present
	if extendsVal, exists := heryContent["_extends"]; exists {
		if _, ok := extendsVal.(string); !ok {
			return fmt.Errorf("_extends must be a string URI")
		}
	}

	// Validate _requires if present
	if requiresVal, exists := heryContent["_requires"]; exists {
		requiresList, ok := requiresVal.([]any)
		if !ok {
			return fmt.Errorf("_requires must be an array of strings")
		}
		for _, item := range requiresList {
			ref, ok := item.(string)
			if !ok {
				return fmt.Errorf("_requires items must be strings")
			}
			if strings.Contains(ref, "../") {
				return fmt.Errorf("_requires references cannot contain '../': %s", ref)
			}
		}
	}

	return nil
}

// Entity validates the YAML content against the JSON schema
func (s *validator) Entity(schema *jsonschema.Schema, heryContent map[string]any) error {
	// 1. Validate _type is present
	typeVal, ok := heryContent["_type"].(string)
	if !ok || typeVal == "" {
		return fmt.Errorf("_type is required")
	}

	// 2. Validate the hery file content with the loaded schema
	if err := schema.Validate(heryContent); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	return nil
}

// EntityUri validates the module path for go get
//
// A entity URI cannot contain the usual URL elements.
func (s *validator) EntityUri(entityUrl string) bool {
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
