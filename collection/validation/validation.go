package validation

import (
	"fmt"
	"regexp"
)

// Collection
func Collection(name string) error {
	if name == "" {
		return fmt.Errorf("collection name cannot be empty")
	}

	// Example: only allow alphanumeric characters and dashes
	re := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	if !re.MatchString(name) {
		return fmt.Errorf("invalid collection name: %s. It should only contain alphanumeric characters and dashes", name)
	}

	return nil
}
