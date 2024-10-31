package storage

import "testing"

func Test_NewStorageService(t *testing.T) {
	storageService := NewStorageService()

	// Check that the returned value is not nil
	if storageService == nil {
		t.Fatal("Expected NewStorageService to return a non-nil value")
	}

	// Check that the returned value is of the correct type
	// TODO: Why was it here in the first place?
	/*if storageService == nil {
		t.Fatalf("Expected NewStorageService to return a non-nil value, but got nil")
	}*/

	// Optionally, check that the fields in AbsPaths are initialized to their zero values
	if storageService.Storage != "" ||
		storageService.Catalog != "" ||
		storageService.Collection != "" ||
		storageService.Entities != "" ||
		storageService.Cache != "" {
		t.Error("Expected all fields of AbsPaths to be initialized to empty strings")
	}
}
