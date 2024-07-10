package validation

import (
	"github.com/AmadlaOrg/hery/collection"
	"log"
	"regexp"
)

// CollectionName validates that the collection name follows the format
func CollectionName(collectionName string) bool {
	matched, err := regexp.MatchString(collection.Match, collectionName)
	if err != nil {
		log.Println("Error validating collection name: ", err)
	}
	return matched
}
