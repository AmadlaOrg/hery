package validation

import (
	"fmt"
	schemaPkg "github.com/AmadlaOrg/hery/entity/schema"
	schemaValidationPkg "github.com/AmadlaOrg/hery/entity/schema/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"path/filepath"
	"strings"
	"unicode"
)

// IValidation used by mockery
type IValidation interface {
	Entity(entityPath, collectionName, entityUri string, heryContent map[string]any) error
	EntityUri(entityUrl string) bool
}

// SValidation used by mockery
type SValidation struct {
	Version           version.IVersion
	VersionValidation versionValidationPkg.IValidation
	Schema            schemaPkg.ISchema
	SchemaValidation  schemaValidationPkg.IValidation
}

// Entity validates the YAML content against the JSON schema
// TODO: Make sure that YAML standard is valid first
// TODO: Since JSON-Schema cannot merge by-it-self the schemas you will need to add code for that
// TODO: Make sure it validates properly with both the based schema found in `.schema` and the entity's own `schema.json`
func (s *SValidation) Entity(entityPath, collectionName, entityUri string, heryContent map[string]any) error {

	// TODO: Is the root of the entity heryContent is equal to the _entity (don't forget version)
	// TODO: If it does not match then check the _self (it does not need to have _entity so it can go straight to validation)
	// TODO: If _entity is in _self then check if it matches if not throw error that it is not valid and why

	// TODO: We need to add a unit test to see what happens when a YAML is not valid in the `.hery` content

	// 1. Get the schema of the entity and load the jsonschema
	// TODO: Move the path generator to schema package
	entityConfigDir := fmt.Sprintf(".%s", collectionName)
	schemaPath := filepath.Join(entityPath, entityConfigDir, schemaPkg.EntityJsonSchemaFileName)
	schema, err := s.Schema.Load(schemaPath)
	if err != nil {
		return fmt.Errorf("error loading JSON schema: %w", err)
	}

	// 2. Validate the hery file content with the loaded schema
	if err = schema.Validate(heryContent); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	// 3. Validate JSON-Schema entity `id`
	err = s.SchemaValidation.Id(schema.ID, collectionName, entityUri)
	if err != nil {
		return err
	}

	// 4. This step is for when the entity is valid
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
