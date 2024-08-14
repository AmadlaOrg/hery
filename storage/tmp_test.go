package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReplaceWithTempDir(t *testing.T) {
	// Test input
	paths := &AbsPaths{
		Storage:    "/some/path/to/.hery/",
		Collection: "/some/path/to/.hery/collection/",
		Entities:   "/some/path/to/.hery/collection/collectionName/entity",
		Cache:      "/some/path/to/.hery/collection/collectionName/collectionName.cache",
	}

	collectionName := "collectionName"

	// Call the function
	newPaths, err := ReplaceWithTempDir(paths, collectionName)
	if err != nil {
		t.Fatalf("ReplaceWithTempDir returned an error: %v", err)
	}

	// Check if the temp directory was correctly created
	tempDir := filepath.Dir(newPaths.Storage)
	if !strings.HasPrefix(tempDir, os.TempDir()) {
		t.Errorf("Expected new paths to start with temp directory, got %s", tempDir)
	}

	// Validate that each path is correctly replaced
	expectedSuffix := filepath.Join(".collectionName", "storage")
	if !strings.HasSuffix(newPaths.Storage, expectedSuffix) {
		t.Errorf("Expected storage path to end with %s, got %s", expectedSuffix, newPaths.Storage)
	}

	expectedSuffix = filepath.Join(".collectionName", "collection")
	if !strings.HasSuffix(newPaths.Collection, expectedSuffix) {
		t.Errorf("Expected collection path to end with %s, got %s", expectedSuffix, newPaths.Collection)
	}

	expectedSuffix = filepath.Join(".collectionName", "entities")
	if !strings.HasSuffix(newPaths.Entities, expectedSuffix) {
		t.Errorf("Expected entities path to end with %s, got %s", expectedSuffix, newPaths.Entities)
	}

	expectedSuffix = filepath.Join(".collectionName", "cache")
	if !strings.HasSuffix(newPaths.Cache, expectedSuffix) {
		t.Errorf("Expected cache path to end with %s, got %s", expectedSuffix, newPaths.Cache)
	}

	// Check for invalid collectionName case
	invalidPaths := &AbsPaths{
		Storage:    "/some/path/to/",
		Collection: "/some/path/to/collection",
		Entities:   "/some/path/to/entities",
		Cache:      "/some/path/to/cache",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when collectionName is not found in paths, but got none")
		}
	}()

	// This should cause a panic due to invalid paths
	_, _ = ReplaceWithTempDir(invalidPaths, collectionName)
}
