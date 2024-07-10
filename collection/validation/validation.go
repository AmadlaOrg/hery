package validation

import (
	"github.com/AmadlaOrg/hery/collection"
	"log"
	"regexp"
)

// Name validates that the collection name follows the format
func Name(collectionName string) bool {
	matched, err := regexp.MatchString(collection.Match, collectionName)
	if err != nil {
		log.Fatalln("Error validating collection name: ", err)
	}
	return matched
}
