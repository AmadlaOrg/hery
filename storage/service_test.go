package storage

import "testing"

func Test_NewStorageService(t *testing.T) {
	storageService := New()

	if storageService == nil {
		t.Fatal("Expected New to return a non-nil value")
	}

	if storageService.Storage != "" ||
		storageService.Entities != "" {
		t.Error("Expected all fields of AbsPaths to be initialized to empty strings")
	}
}
