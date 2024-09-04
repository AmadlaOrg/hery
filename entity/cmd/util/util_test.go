package util

import (
	"bytes"
	"fmt"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestConcoct_Success(t *testing.T) {
	expectedPaths := storage.AbsPaths{
		Storage:    "",
		Catalog:    "",
		Collection: "",
		Entities:   "",
		Cache:      "",
	}

	getCollectionFlag = func() (string, error) {
		return "testCollection", nil
	}

	mockStorage := storage.MockStorage{}
	mockStorage.EXPECT().Paths(mock.Anything).Return(&expectedPaths, nil)

	handlerCalled := false
	handler := func(collectionName string, paths *storage.AbsPaths, args []string) {
		handlerCalled = true
		require.Equal(t, "testCollection", collectionName)
		require.Equal(t, &expectedPaths, paths)
		require.Empty(t, args) // Assuming no additional args are passed
	}

	mockSUtil := SUtil{
		newStorageService: &mockStorage,
	}

	mockSUtil.Concoct(&cobra.Command{}, []string{}, handler)

	require.True(t, handlerCalled)
	mockStorage.AssertExpectations(t)
}

func TestConcoct_StoragePathsError(t *testing.T) {
	// Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(nil) // Reset log output after test

	// Mock getCollectionFlag to return a valid collection name
	getCollectionFlag = func() (string, error) {
		return "testCollection", nil
	}

	// Mock storage to return an error when Paths is called
	mockStorage := storage.MockStorage{}
	mockStorage.EXPECT().Paths("testCollection").Return(nil, fmt.Errorf("mock error from Paths"))

	mockSUtil := SUtil{
		newStorageService: &mockStorage,
	}

	handlerCalled := false
	handler := func(collectionName string, paths *storage.AbsPaths, args []string) {
		handlerCalled = true
	}

	// Concoct should log the error and not call the handler
	mockSUtil.Concoct(&cobra.Command{}, []string{}, handler)

	// Assert log output
	require.Contains(t, logBuffer.String(), "Error getting paths: mock error from Paths")

	// Ensure handler was not called
	require.False(t, handlerCalled)
	mockStorage.AssertExpectations(t)
}

func TestConcoct_GetCollectionFlagError(t *testing.T) {
	// Mock getCollectionFlag to return an error
	/*getCollectionFlag = func() (string, error) {
		return "", errors.New("getCollectionFlag error")
	}

	// Mock storage service won't be called due to error
	newStorageService = func() storage.IStorage {
		return new(MockStorageService)
	}

	// Capture the log output
	logOutput := captureLogOutput(func() {
		Concoct(&cobra.Command{}, []string{}, nil)
	})

	// Assert log output contains the error
	assert.Contains(t, logOutput, "getCollectionFlag error")

	// Clean up
	getCollectionFlag = collectionPkgCmd.GetCollectionFlag
	newStorageService = storage.NewStorageService*/
}

func TestConcoct_PathsError(t *testing.T) {
	// Replace the function variables with mocks
	/*getCollectionFlag = mockGetCollectionFlag
	mockStorageService := new(MockStorageService)
	newStorageService = func() storage.IStorage {
		return mockStorageService
	}

	// Prepare mock responses
	mockStorageService.On("Paths", "mockCollection").Return(nil, errors.New("paths error"))

	// Capture the log output
	logOutput := captureLogOutput(func() {
		Concoct(&cobra.Command{}, []string{}, nil)
	})

	// Assert log output contains the error
	assert.Contains(t, logOutput, "Error getting paths: paths error")

	// Clean up
	getCollectionFlag = collectionPkgCmd.GetCollectionFlag
	newStorageService = storage.NewStorageService*/
}

/*func TestConcoct_GetCollectionFlagError(t *testing.T) {
	// Mock dependencies
	mockCollectionCmd := new(MockCollectionCmd)

	// Override the actual functions with mocks
	getCollectionFlag = mockCollectionCmd.GetCollectionFlag

	// Set up expectations
	mockCollectionCmd.On("GetCollectionFlag").Return("", errors.New("collection flag error"))

	// Capture log output
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	// Run the test
	handler := func(collectionName string, paths *storage.AbsPaths, args []string) {
		t.Fail() // Handler should not be called
	}
	cobraCmd := &cobra.Command{}
	cobraCmd.Run = func(cmd *cobra.Command, args []string) {
		concoct(cmd, args, handler)
	}
	cobraCmd.Run(cobraCmd, nil)

	assert.Contains(t, logOutput.String(), "collection flag error")

	// Verify expectations
	mockCollectionCmd.AssertExpectations(t)
}*/

/*func TestConcoct_PathsError(t *testing.T) {
	// Mock dependencies
	mockStorageService := new(MockStorageService)
	mockCollectionCmd := new(MockCollectionCmd)

	// Override the actual functions with mocks
	getCollectionFlag = mockCollectionCmd.GetCollectionFlag
	newStorageService = func() *storage.AbsPaths {
		return mockStorageService
	}

	// Set up test data
	expectedCollectionName := "testCollection"

	// Set up expectations
	mockCollectionCmd.On("GetCollectionFlag").Return(expectedCollectionName, nil)
	mockStorageService.On("Paths", expectedCollectionName).Return(nil, errors.New("paths error"))

	// Capture log output
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	// Run the test
	handler := func(collectionName string, paths *storage.AbsPaths, args []string) {
		t.Fail() // Handler should not be called
	}
	cobraCmd := &cobra.Command{}
	cobraCmd.Run = func(cmd *cobra.Command, args []string) {
		concoct(cmd, args, handler)
	}
	cobraCmd.Run(cobraCmd, nil)

	assert.Contains(t, logOutput.String(), "Error getting paths: paths error")

	// Verify expectations
	mockCollectionCmd.AssertExpectations(t)
	mockStorageService.AssertExpectations(t)
}*/
