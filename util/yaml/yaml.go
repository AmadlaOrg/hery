package yaml

import (
	"fmt"
	"github.com/AmadlaOrg/hery/util/file"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Yaml interface {
	Read(path string, fileName string) (map[string]interface{}, error)
}

// Read makes it easy to read any yaml file with any of the two extensions: yml or yaml
func Read(path string, fileName string) (map[string]interface{}, error) {
	fullFileNameYml := fmt.Sprintf("%s.yml", fileName)
	fullFileNameYaml := fmt.Sprintf("%s.yaml", fileName)

	var fullPath string
	YmlPath := filepath.Join(path, fullFileNameYml)
	YamlPath := filepath.Join(path, fullFileNameYaml)

	YmlPathExists := file.Exists(YmlPath)
	YamlPathExists := file.Exists(YamlPath)

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
	if err != nil {
		return nil, err
	}

	return current, nil
}
