package validation

import (
	"fmt"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/AmadlaOrg/hery/util/file"
)

// AllExist validates all the paths
func AllExist(absPaths storage.AbsPaths) error {
	if !file.Exists(absPaths.Storage) {
		return fmt.Errorf("storage %s does not exist", absPaths.Storage)
	}
	if !file.Exists(absPaths.Collection) {
		return fmt.Errorf("collection %s does not exist", absPaths.Storage)
	}
	if !file.Exists(absPaths.Entities) {
		return fmt.Errorf("entities %s does not exist", absPaths.Storage)
	}
	if !file.Exists(absPaths.Cache) {
		return fmt.Errorf("cache %s does not exist", absPaths.Storage)
	}
	return nil
}
