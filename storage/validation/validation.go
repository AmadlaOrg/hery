package validation

import (
	"fmt"
	"github.com/AmadlaOrg/LibraryUtils/file"
	"github.com/AmadlaOrg/hery/storage"
)

var fileExists = file.Exists

// AllExist validates all the paths
func AllExist(absPaths storage.AbsPaths) error {
	if !fileExists(absPaths.Storage) {
		return fmt.Errorf("storage %s does not exist", absPaths.Storage)
	}
	if !fileExists(absPaths.Entities) {
		return fmt.Errorf("entities %s does not exist", absPaths.Entities)
	}
	return nil
}
