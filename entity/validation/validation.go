package validation

import (
	"github.com/santhosh-tekuri/jsonschema"
	"strings"
	"unicode"
)

// Schema
func Schema() *jsonschema.Schema {
	return nil
}

// Entity
func Entity() error {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7
	return nil
}

// EntityUrl validates the module path for go get
func EntityUrl(path string) bool {
	if strings.Contains(path, "://") {
		return false
	}
	for _, r := range path {
		if unicode.IsSpace(r) || r == ':' || r == '?' || r == '&' || r == '=' || r == '#' {
			return false
		}
	}
	return true
}
