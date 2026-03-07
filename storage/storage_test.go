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

func TestTmpMain(t *testing.T) {
	setup()
	d := &AbsPaths{}
	mainPath, err := d.TmpMain()
	assert.NoError(t, err)
	assert.Equal(t, "/tmp/mock/hery", mainPath)
}

func TestMakePaths(t *testing.T) {
	setup()
	d := &AbsPaths{}
	paths := AbsPaths{
		Storage:  "/mock/path/storage",
		Entities: "/mock/path/entities",
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
		Storage:  "/mock/path/storage",
		Entities: "/mock/path/entities",
	}
	err := d.MakePaths(paths)
	assert.Error(t, err)
	assert.Equal(t, "mock error", err.Error())
}

// TestMainPathUsingEnvVar tests the Main function using an environment variable
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

// TestMainPathUsingDefault tests the Main function using the default path
func TestMainPathUsingDefault(t *testing.T) {
	mockStorage := NewMockStorage(t)

	var expectedPath string
	switch runtime.GOOS {
	case "windows":
		expectedPath = filepath.Join(os.Getenv("APPDATA"), "hery")
	default:
		homeDir, _ := os.UserHomeDir()
		expectedPath = filepath.Join(homeDir, ".cache", "hery")
	}

	mockStorage.On("Main").Return(expectedPath, nil)

	actualPath, err := mockStorage.Main()

	assert.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)

	mockStorage.AssertExpectations(t)
}

// TestEntityPath tests the EntityPath function
func TestEntityPath(t *testing.T) {
	mockStorage := NewMockStorage(t)

	entitiesPath := "/mock/storage/entity"
	entityRelativePath := "some/entity.yaml"
	expectedPath := filepath.Join(entitiesPath, entityRelativePath)

	mockStorage.On("EntityPath", entitiesPath, entityRelativePath).Return(expectedPath)

	actualPath := mockStorage.EntityPath(entitiesPath, entityRelativePath)

	assert.Equal(t, expectedPath, actualPath)

	mockStorage.AssertExpectations(t)
}

// TestMainError tests error case
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
