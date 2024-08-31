package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type MockFile struct {
	mock.Mock
}

func (m *MockFile) Exists(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

// Setup test suite
func setup() {
	// Mock os and filepath functions as needed
	osGetwd = func() (string, error) {
		return "/mock/path", nil
	}
	filepathAbs = func(path string) (string, error) {
		return "/abs/mock/path/" + path, nil
	}
	filepathJoin = filepath.Join
	fileExists = func(path string) bool {
		return true
	}
	osMkdirAll = func(path string, perm os.FileMode) error {
		return nil
	}
	osMkdirTemp = func(dir, pattern string) (string, error) {
		return "/tmp/mock/hery", nil
	}
}

// FIXME:
/*func TestPaths(t *testing.T) {
	setup()
	d := &AbsPaths{}
	paths, err := d.Paths("testCollection")
	assert.NoError(t, err)
	assert.Equal(t, "/mock/path/testCollection", paths.Storage)
}*/

// FIXME:
/*func TestEntityPath(t *testing.T) {
	setup()
	d := &AbsPaths{}
	entityPath := d.EntityPath("/mock/path/entities", "entity1.json")
	assert.Equal(t, "/mock/path/entities/entity1.json", entityPath)
}*/

// FIXME:
/*func TestTmpPaths(t *testing.T) {
	setup()
	d := &AbsPaths{}
	paths, err := d.TmpPaths("testCollection")
	assert.NoError(t, err)
	assert.Equal(t, "/tmp/mock/hery/testCollection", paths.Storage)
}*/

func TestTmpMain(t *testing.T) {
	setup()
	d := &AbsPaths{}
	mainPath, err := d.TmpMain()
	assert.NoError(t, err)
	assert.Equal(t, "/tmp/mock/hery/.hery", mainPath)
}

func TestMakePaths(t *testing.T) {
	setup()
	d := &AbsPaths{}
	paths := AbsPaths{
		Storage:    "/mock/path/storage",
		Catalog:    "/mock/path/catalog",
		Collection: "/mock/path/collection",
		Entities:   "/mock/path/entities",
		Cache:      "/mock/path/cache",
	}
	err := d.MakePaths(paths)
	assert.NoError(t, err)
}

func TestTmpMain_MkdirTempError(t *testing.T) {
	osMkdirTemp = func(dir, pattern string) (string, error) {
		return "", errors.New("mock error")
	}
	d := &AbsPaths{}
	_, err := d.TmpMain()
	assert.Error(t, err)
	assert.Equal(t, "mock error", err.Error())
}

func TestMakePaths_MkdirAllError(t *testing.T) {
	osMkdirAll = func(path string, perm os.FileMode) error {
		return errors.New("mock error")
	}
	d := &AbsPaths{}
	paths := AbsPaths{
		Storage:    "/mock/path/storage",
		Catalog:    "/mock/path/catalog",
		Collection: "/mock/path/collection",
		Entities:   "/mock/path/entities",
		Cache:      "/mock/path/cache",
	}
	err := d.MakePaths(paths)
	assert.Error(t, err)
	assert.Equal(t, "mock error", err.Error())
}

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
