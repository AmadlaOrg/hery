package validation

import (
	"encoding/json"
	"fmt"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/util/file"
	schemaPkg "github.com/AmadlaOrg/hery/util/json/schema"
	"github.com/santhosh-tekuri/jsonschema"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"unicode"
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
func Entity(entityPath string) error {
	schemaPath := filepath.Join(entityPath, ".amadla", "schema.json")
	schema, err := schemaPkg.Load(schemaPath)
	if err != nil {
		return fmt.Errorf("error loading JSON schema: %w", err)
	}

	yamlFilePath := filepath.Join(entityPath, "entity.yaml") // Assume the YAML file name
	if !file.Exists(yamlFilePath) {
		yamlFilePath = filepath.Join(entityPath, "entity.yml")
	}
	if !file.Exists(yamlFilePath) {
		return fmt.Errorf("YAML file not found in entity: %s", entityPath)
	}

	yamlContent, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %w", err)
	}

	var yamlData map[string]interface{}
	if err := yaml.Unmarshal(yamlContent, &yamlData); err != nil {
		return fmt.Errorf("error unmarshalling YAML content: %w", err)
	}

	jsonData, err := json.Marshal(yamlData)
	if err != nil {
		return fmt.Errorf("error marshalling YAML content to JSON: %w", err)
	}

	var jsonDataInterface interface{}
	if err := json.Unmarshal(jsonData, &jsonDataInterface); err != nil {
		return fmt.Errorf("error unmarshalling JSON content: %w", err)
	}

	if err := schema.ValidateInterface(jsonDataInterface); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	return nil
}

// EntityUrl validates the module path for go get
//
// -------------------------------------------------------------------------------
//
// --- Params ---
// entityUrl string: the entity url (with version or not)
// checkVersionExist bool: to indicate to check the version exist or not
//
// --- Return ---
// bool: is valid true or not
// string: valid version
// error: error message
//
// -------------------------------------------------------------------------------
func EntityUrl(entityUrl string, checkVersionExist bool) (bool, string, error) {
	if strings.Contains(entityUrl, "://") {
		return false, "", nil
	}
	for _, r := range entityUrl {
		if unicode.IsSpace(r) || r == ':' || r == '?' || r == '&' || r == '=' || r == '#' {
			return false, "", nil
		}
	}
	if strings.Contains(entityUrl, "@") {
		ver, err := version.Extract(entityUrl)
		if err != nil {
			return false, "", fmt.Errorf("error extracting version: %v", err)
		}
		if !versionValidationPkg.Format(ver) {
			return false, "", nil
		}
		// TODO: Check with git if the version actually exists
	}
	return true, "", nil
}
