package env

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create the env subdirectory within the temporary directory
	envDir := filepath.Join(tmpDir, "env")
	err := os.Mkdir(envDir, 0755)
	assert.NoError(t, err)

	// Create a temporary types.go file in the env subdirectory
	typesGoContent := `
	package env

	const (
		Const1 = "Value1"
		Const2 = "Value2"
		Const3 = "Value3"
	)
	`
	typesGoPath := filepath.Join(envDir, "types.go")
	err = os.WriteFile(typesGoPath, []byte(typesGoContent), 0644)
	assert.NoError(t, err)

	// Override the filepath to point to the temporary directory
	originalDir, err := os.Getwd()
	assert.NoError(t, err)
	defer func() {
		err := os.Chdir(originalDir)
		if err != nil {
			t.Fatalf("Failed to change back to the original directory: %v", err)
		}
	}()
	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	// Call the List function and check the result
	expectedConstants := []string{"Value1", "Value2", "Value3"}
	actualConstants, err := List()
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedConstants, actualConstants)

	// Test error handling for absolute path failure
	t.Run("absolute path failure", func(t *testing.T) {
		// Save the original filepathAbs function and restore it after the test
		originalFilepathAbs := filepathAbs
		defer func() { filepathAbs = originalFilepathAbs }()

		// Replace filepathAbs with a function that returns an error
		filepathAbs = func(_ string) (string, error) {
			return "", fmt.Errorf("simulated failure in filepath.Abs")
		}

		_, err := List()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get the absolute path of the current directory")
	})

	// Test error handling for parsing file failure
	t.Run("parse file failure", func(t *testing.T) {
		// Create an invalid types.go file
		invalidTypesGoContent := `
		package env

		const (
			Const1 = "Value1
		)
		`
		err := os.WriteFile(typesGoPath, []byte(invalidTypesGoContent), 0644)
		assert.NoError(t, err)

		_, err = List()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse file")
	})

	// Test handling for non-constant declarations
	t.Run("non-constant declarations", func(t *testing.T) {
		nonConstantTypesGoContent := `
		package env

		var (
			Var1 = "Value1"
		)
		`
		err := os.WriteFile(typesGoPath, []byte(nonConstantTypesGoContent), 0644)
		assert.NoError(t, err)

		constants, err := List()
		assert.NoError(t, err)
		assert.Empty(t, constants)
	})

	// Test handling for non-string constants
	t.Run("non-string constants", func(t *testing.T) {
		nonStringTypesGoContent := `
		package env

		const (
			IntConst = 1
			FloatConst = 1.1
			BoolConst = true
		)
		`
		err := os.WriteFile(typesGoPath, []byte(nonStringTypesGoContent), 0644)
		assert.NoError(t, err)

		constants, err := List()
		assert.NoError(t, err)
		assert.Empty(t, constants)
	})

	// Test handling for type declarations
	t.Run("type declarations", func(t *testing.T) {
		typeDeclarationsContent := `
		package env

		type MyType struct {
			Field1 string
			Field2 int
		}
		`
		err := os.WriteFile(typesGoPath, []byte(typeDeclarationsContent), 0644)
		assert.NoError(t, err)

		constants, err := List()
		assert.NoError(t, err)
		assert.Empty(t, constants)
	})
}
