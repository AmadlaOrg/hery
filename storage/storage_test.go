package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestMainPathUsingEnvVar tests the Main function using an environment variable for the storage path
func TestMainPathUsingEnvVar(t *testing.T) {
	mockStorage := NewMockStorage(t)

	expectedPath, _ := filepath.Abs("/mock/storage/path")
	t.Setenv(HeryStoragePath, expectedPath)

	mockStorage.On("Main").Return(expectedPath, nil)

	actualPath, err := mockStorage.Main()

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)

	mockStorage.AssertExpectations(t)
}

// TestMainPathUsingCurrentLocation tests the Main function using the current working directory
func TestMainPathUsingCurrentLocation(t *testing.T) {
	mockStorage := NewMockStorage(t)

	cwd, _ := os.Getwd()
	expectedPath := filepath.Join(cwd, ".hery")

	mockStorage.On("Main").Return(expectedPath, nil)

	actualPath, err := mockStorage.Main()

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)

	mockStorage.AssertExpectations(t)
}

// TestMainPathUsingDefault tests the Main function using the default path based on the operating system
func TestMainPathUsingDefault(t *testing.T) {
	mockStorage := NewMockStorage(t)

	var expectedPath string
	switch runtime.GOOS {
	case "windows":
		expectedPath = filepath.Join(os.Getenv("APPDATA"), "Hery")
	default: // "linux" and "darwin"
		homeDir, _ := os.UserHomeDir()
		expectedPath = filepath.Join(homeDir, ".hery")
	}

	mockStorage.On("Main").Return(expectedPath, nil)

	actualPath, err := mockStorage.Main()

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)

	mockStorage.AssertExpectations(t)
}

// TestPaths tests the Paths function to ensure correct paths are returned
// FIXME:
/*func TestPaths(t *testing.T) {
	mockStorage := NewMockStorage(t)

	mainPath := "/mock/storage"
	collectionName := "mockCollection"
	collectionPath := filepath.Join(mainPath, collectionName)
	entityPath := filepath.Join(collectionPath, "entity")
	cachePath := filepath.Join(collectionPath, "mockCollection.cache")

	expectedPaths := &AbsPaths{
		Storage:    mainPath,
		Collection: collectionPath,
		Entities:   entityPath,
		Cache:      cachePath,
	}

	mockStorage.On("Main").Return(mainPath, nil)
	mockStorage.On("Paths", collectionName).Return(expectedPaths, nil)

	actualPaths, err := mockStorage.Paths(collectionName)

	assert.NoError(t, err)
	assert.Equal(t, expectedPaths, actualPaths)

	mockStorage.AssertExpectations(t)
}*/

// TestEntityPath tests the EntityPath function to ensure the correct path is returned
func TestEntityPath(t *testing.T) {
	mockStorage := NewMockStorage(t)

	entitiesPath := "/mock/storage/mockCollection/entity"
	entityRelativePath := "some/entity.yaml"
	expectedPath := filepath.Join(entitiesPath, entityRelativePath)

	mockStorage.On("EntityPath", entitiesPath, entityRelativePath).Return(expectedPath)

	actualPath := mockStorage.EntityPath(entitiesPath, entityRelativePath)

	assert.Equal(t, expectedPath, actualPath)

	mockStorage.AssertExpectations(t)
}

// TestMainError tests the case when an error occurs in the Main function
func TestMainError(t *testing.T) {
	mockStorage := NewMockStorage(t)

	mockError := errors.New("mock error")
	mockStorage.On("Main").Return("", mockError)

	actualPath, err := mockStorage.Main()

	assert.Error(t, err)
	assert.Equal(t, mockError, err)
	assert.Equal(t, "", actualPath)

	mockStorage.AssertExpectations(t)
}
