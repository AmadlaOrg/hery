package heryext

import (
	"fmt"
	"github.com/AmadlaOrg/hery/util/file"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type IHeryExt interface {
	Read(path, collectionName string) (map[string]any, error)
}

type SHeryExt struct{}

// Read makes it easy to read any yaml file with any of the two extensions: yml or yaml
func (s *SHeryExt) Read(path, collectionName string) (map[string]any, error) {
	heryFileName := fmt.Sprintf("%s.hery", collectionName)
	heryPath := filepath.Join(path, heryFileName)

	if !file.Exists(heryPath) {
		return nil, fmt.Errorf("%s does not exist", heryPath)
	}

	content, err := os.ReadFile(heryPath)
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
