package validation

import (
	"fmt"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/AmadlaOrg/hery/util/file"
)

var fileExists = file.Exists

// AllExist validates all the paths
func AllExist(absPaths storage.AbsPaths) error {
	if !fileExists(absPaths.Storage) {
		return fmt.Errorf("storage %s does not exist", absPaths.Storage)
	}
	if !fileExists(absPaths.Catalog) {
		return fmt.Errorf("catalog %s does not exist", absPaths.Catalog)
	}
	if !fileExists(absPaths.Collection) {
		return fmt.Errorf("collection %s does not exist", absPaths.Storage)
	}
	if !fileExists(absPaths.Entities) {
		return fmt.Errorf("entities %s does not exist", absPaths.Storage)
	}
	if !fileExists(absPaths.Cache) {
		return fmt.Errorf("cache %s does not exist", absPaths.Storage)
	}
	return nil
}
