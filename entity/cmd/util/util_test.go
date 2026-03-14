package util

import (
	"bytes"
	"fmt"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestConcoct_Success(t *testing.T) {
	expectedPaths := storage.AbsPaths{
		Storage:  "",
		Entities: "",
		Cache:    "",
	}

	mockStorage := storage.MockStorage{}
	mockStorage.EXPECT().Paths().Return(&expectedPaths, nil)

	handlerCalled := false
	handler := func(paths *storage.AbsPaths, args []string) {
		handlerCalled = true
		require.Equal(t, &expectedPaths, paths)
		require.Empty(t, args)
	}

	mockSUtil := utilImpl{
		New: &mockStorage,
	}

	err := mockSUtil.Concoct(&cobra.Command{}, []string{}, handler)
	if err != nil {
		t.Fatal(err)
	}

	require.True(t, handlerCalled)
	mockStorage.AssertExpectations(t)
}

func TestConcoct_StoragePathsError(t *testing.T) {
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(nil)

	mockStorage := storage.MockStorage{}
	mockStorage.EXPECT().Paths().Return(nil, fmt.Errorf("mock error from Paths"))

	mockSUtil := utilImpl{
		New: &mockStorage,
	}

	handlerCalled := false
	handler := func(paths *storage.AbsPaths, args []string) {
		handlerCalled = true
	}

	err := mockSUtil.Concoct(&cobra.Command{}, []string{}, handler)

	require.Error(t, err)
	require.Contains(t, err.Error(), "mock error from Paths")

	require.False(t, handlerCalled)
	mockStorage.AssertExpectations(t)
}
