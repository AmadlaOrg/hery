package validation

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
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
	Entity(collectionName string, schema *jsonschema.Schema, heryContent entity.NotFormatedContent) error
	EntityUri(entityUrl string) bool
}

// SValidation used by mockery
type SValidation struct {
	Version           version.IVersion
	VersionValidation versionValidationPkg.IValidation
	Schema            schemaPkg.ISchema
	SchemaValidation  schemaValidationPkg.IValidation
}

// RootEntity
func (s *SValidation) RootEntity(rootSchema, selfSchema *jsonschema.Schema, heryContent map[string]any) error {

	// TODO: Is the root of the entity heryContent is equal to the _entity (don't forget version)
	// TODO: If it does not match then check the _self (it does not need to have _entity so it can go straight to validation)
	// TODO: If _entity is in _self then check if it matches if not throw error that it is not valid and why

	// . This step is for when the root entity with the _self entity are valid is valid
	return nil
}

// Entity validates the YAML content against the JSON schema
func (s *SValidation) Entity(
	collectionName string, schema *jsonschema.Schema, heryContent entity.NotFormatedContent) error {

	// TODO: We need to add a unit test to see what happens when a YAML is not valid in the `.hery` content
	// |-- TODO: Make sure that YAML standard is valid first

	// 1. Get the schema of the entity and load the jsonschema
	// TODO: Move outside so that it is added to type Entity
	/*schema, err := s.Schema.Load(schemaPath)
	if err != nil {
		return fmt.Errorf("error loading JSON schema: %w", err)
	}*/

	// 1. Validate JSON-Schema entity `id`
	// TODO: Move it somewhere else
	err := s.SchemaValidation.Id(schema.ID, collectionName, heryContent["_entity"].(string))
	if err != nil {
		return err
	}

	// 2. Validate the hery file content with the loaded schema
	if err = schema.Validate(heryContent); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	// 3. This step is for when the entity is valid
	return nil
}

func (s *SValidation) Body() {

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
