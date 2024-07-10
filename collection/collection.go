package collection

import (
	"path/filepath"
)

// Path returns the collection absolute path
func Path(collectionName, storagePath string) string {
	return filepath.Join(storagePath, collectionName)
}
