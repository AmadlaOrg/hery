package env

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strings"
)

// List collects string constants from the types.go file in the env package
func List() []string {
	// Get the absolute path of the current directory
	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Fatalf("Failed to get the absolute path of the current directory: %v", err)
	}

	// Construct the full path to types.go
	filename := filepath.Join(dir, "env", "types.go")
	var constants []string

	// Parse the Go source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Failed to parse file %s: %v", filename, err)
	}

	// Inspect the AST and collect string constants
	ast.Inspect(node, func(n ast.Node) bool {
		decl, ok := n.(*ast.GenDecl)
		if !ok || decl.Tok != token.CONST {
			return true
		}

		for _, spec := range decl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for _, value := range valueSpec.Values {
				basicLit, ok := value.(*ast.BasicLit)
				if !ok || basicLit.Kind != token.STRING {
					continue
				}
				// Trim the quotes from the string value
				constants = append(constants, strings.Trim(basicLit.Value, "\""))
			}
		}
		return false
	})

	return constants
}
