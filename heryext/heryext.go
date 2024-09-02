package heryext

import (
	"fmt"
	"github.com/AmadlaOrg/hery/util/file"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type IHeryExt interface {
	Read(path string, fileName string) (map[string]interface{}, error)
}

type SHeryExt struct{}

// Read makes it easy to read any yaml file with any of the two extensions: yml or yaml
func (s *SHeryExt) Read(path string, fileName string) (map[string]interface{}, error) {
	fullFileNameYaml := fmt.Sprintf("%s.hery", fileName)

	var fullPath string
	YamlPath := filepath.Join(path, fullFileNameYaml)

	YamlPathExists := file.Exists(YamlPath)

	if !YamlPathExists {
		return nil, fmt.Errorf("%s does not exist", YamlPath)
	} else {
		fullPath = YamlPath
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	var current map[string]interface{}
	err = yaml.Unmarshal(content, &current)
	if err != nil {
		return nil, err
	}

	return current, nil
}
