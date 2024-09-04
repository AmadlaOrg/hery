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

	err := mockSUtil.Concoct(&cobra.Command{}, []string{}, handler)
	if err != nil {
		t.Fatal(err)
	}

	require.True(t, handlerCalled)
	mockStorage.AssertExpectations(t)
}

func TestConcoct_GetCollectionFlagError(t *testing.T) {
	// Mock getCollectionFlag to return an error
	getCollectionFlag = func() (string, error) {
		return "", fmt.Errorf("mock error from getCollectionFlag")
	}

	mockStorage := storage.MockStorage{}

	mockSUtil := SUtil{
		newStorageService: &mockStorage,
	}

	handlerCalled := false
	handler := func(collectionName string, paths *storage.AbsPaths, args []string) {
		handlerCalled = true
	}

	err := mockSUtil.Concoct(&cobra.Command{}, []string{}, handler)

	// Assert that an error was returned
	require.Error(t, err)
	require.Contains(t, err.Error(), "mock error from getCollectionFlag")

	// Ensure the handler was not called
	require.False(t, handlerCalled)
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
	err := mockSUtil.Concoct(&cobra.Command{}, []string{}, handler)

	// Assert that an error was returned
	require.Error(t, err)
	require.Contains(t, err.Error(), "mock error from Paths")

	// Ensure handler was not called
	require.False(t, handlerCalled)
	mockStorage.AssertExpectations(t)
}
