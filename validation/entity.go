package validation

import (
	"encoding/json"
	"fmt"
	"github.com/AmadlaOrg/hery/util"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Entity validates the YAML content against the JSON schema
func Entity(entityPath string) error {
	schemaPath := filepath.Join(entityPath, ".amadla", "schema.json")
	schema, err := util.LoadJSONSchema(schemaPath)
	if err != nil {
		return fmt.Errorf("error loading JSON schema: %w", err)
	}

	yamlFilePath := filepath.Join(entityPath, "entity.yaml") // Assume the YAML file name
	if !util.FileExists(yamlFilePath) {
		yamlFilePath = filepath.Join(entityPath, "entity.yml")
	}
	if !util.FileExists(yamlFilePath) {
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
