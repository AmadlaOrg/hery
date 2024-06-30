package validation

import (
	"github.com/AmadlaOrg/hery/collection"
	"regexp"
)

// CollectionName validates that the collection name follows the format
func CollectionName(collectionName string) (bool, error) {
	matched, err := regexp.MatchString(collection.Match, collectionName)
	if err != nil {
		return false, err
	}
	return matched, nil
}
