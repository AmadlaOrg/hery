package util

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// FileExists verify that a file or directory exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// ReadYaml makes it easy to read any yaml file with any of the two extensions: yml or yaml
func ReadYaml(path string, fileName string) (map[string]interface{}, error) {
	fullFileNameYml := fmt.Sprintf("%s.yml", fileName)
	fullFileNameYaml := fmt.Sprintf("%s.yaml", fileName)

	var fullPath string
	YmlPath := filepath.Join(path, fullFileNameYml)
	YamlPath := filepath.Join(path, fullFileNameYaml)

	YmlPathExists := FileExists(YmlPath)
	YamlPathExists := FileExists(YamlPath)

	if !YmlPathExists && !YamlPathExists {
		return nil, fmt.Errorf("%s does not exist", YmlPath)
	} else if YmlPathExists && YamlPathExists {
		return nil, fmt.Errorf("both %s, %s exists", YmlPath, YamlPath)
	} else if YmlPathExists {
		fullPath = YmlPath
	} else {
		fullPath = YamlPath
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	var current map[string]interface{}
	err = yaml.Unmarshal(content, &current)

	return current, nil
}
