package validation

import (
	"fmt"
	schemaPkg "github.com/AmadlaOrg/hery/entity/schema"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// IValidation
type IValidation interface {
	Entity(entityPath, collectionName string, heryContent map[string]any) error
	EntityUri(entityUrl string) bool
}

// SValidation
type SValidation struct {
	Version           version.IVersion
	VersionValidation versionValidationPkg.IValidation
	Schema            schemaPkg.ISchema
}

// Entity validates the YAML content against the JSON schema
// TODO: Make sure that YAML standard is valid first
// TODO: Since JSON-Schema cannot merge by-it-self the schemas you will need to add code for that
// TODO: Make sure it validates properly with both the based schema found in `.schema` and the entity's own `schema.json`
func (s *SValidation) Entity(entityPath, collectionName string, heryContent map[string]any) error {
	// 1. Get the schema of the entity and load the jsonschema
	schemaPath := filepath.Join(entityPath, fmt.Sprintf(".%s", collectionName), "schema.json")
	schema, err := s.Schema.Load(schemaPath)
	if err != nil {
		return fmt.Errorf("error loading JSON schema: %w", err)
	}

	// 2. Validate the hery file content with the loaded schema
	if err = schema.Validate(heryContent); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	// TODO: Put it in a function (it might be used in other context)
	// TODO: Add tests
	if schema.ID == "" {
		return fmt.Errorf("schema validation failed: no ID found in schema")
	}

	idUrnRegex := regexp.MustCompile(schemaPkg.EntityJsonSchemaIdURN)
	if !idUrnRegex.MatchString(schema.ID) {
		return fmt.Errorf("schema validation failed: invalid ID found in schema")
	}

	// TODO: Check for :hery: as first after urn:
	// TODO: Check for collection name
	// TODO: Check for Entity name
	// TODO: Check for version format and if it exist (exist less important)

	return nil
}

// EntityUri validates the module path for go get
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
