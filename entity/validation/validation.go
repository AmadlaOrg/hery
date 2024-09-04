package validation

import (
	"encoding/json"
	"fmt"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	schemaPkg "github.com/AmadlaOrg/hery/heryext/schema"
	"github.com/AmadlaOrg/hery/util/file"
	"github.com/santhosh-tekuri/jsonschema"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// IValidation
type IValidation interface {
	Entity(collectionName, entityPath string) error
	EntityUri(entityUrl string) bool
}

// SValidation
type SValidation struct {
	Version           version.IVersion
	VersionValidation versionValidationPkg.IValidation
	Schema            schemaPkg.ISchema
}

var (
	osReadFile    = os.ReadFile
	jsonMarshal   = json.Marshal
	jsonUnmarshal = json.Unmarshal
)

// Schema
func Schema() *jsonschema.Schema {
	return nil
}

// Entity
/*func Entity() error {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7
	return nil
}*/

// Entity validates the YAML content against the JSON schema
func (s *SValidation) Entity(collectionName, entityPath string) error {
	schemaPath := filepath.Join(entityPath, fmt.Sprintf(".%s", collectionName), "schema.json")
	schema, err := s.Schema.Load(schemaPath)
	if err != nil {
		return fmt.Errorf("error loading JSON schema: %w", err)
	}

	yamlFilePath := filepath.Join(entityPath, fmt.Sprintf("%s.hery", collectionName)) // Assume the YAML file name
	/*if !file.Exists(yamlFilePath) {
		yamlFilePath = filepath.Join(entityPath, fmt.Sprintf("%s.yml", collectionName))
	}*/
	if !file.Exists(yamlFilePath) {
		return fmt.Errorf("HERY file not found in entity: %s", entityPath)
	}

	yamlContent, err := osReadFile(yamlFilePath)
	if err != nil {
		return fmt.Errorf("error reading HERY file: %w", err)
	}

	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(yamlContent, &yamlData); err != nil {
		return fmt.Errorf("error unmarshalling HERY content: %w", err)
	}

	jsonData, err := jsonMarshal(yamlData)
	if err != nil {
		return fmt.Errorf("error marshalling HERY content to JSON: %w", err)
	}

	var jsonDataInterface interface{}
	if err := jsonUnmarshal(jsonData, &jsonDataInterface); err != nil {
		return fmt.Errorf("error unmarshalling JSON content: %w", err)
	}

	if err := schema.ValidateInterface(jsonDataInterface); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

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
