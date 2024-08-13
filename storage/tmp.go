package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ReplaceWithTempDir
func ReplaceWithTempDir(paths *AbsPaths, collectionName string) (*AbsPaths, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "tempdir_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Replace the prefix before .<collectionName> with the tempDir path
	replacePath := func(originalPath string) string {
		// Find the part of the path that starts with .<collectionName>
		collectionIndex := strings.Index(originalPath, fmt.Sprintf(".%s", collectionName))
		if collectionIndex == -1 {
			log.Fatalf("collectionName '%s' not found in path: %s", collectionName, originalPath)
		}
		// Join the tempDir with the part after .<collectionName>
		return filepath.Join(tempDir, originalPath[collectionIndex:])
	}

	// Return new paths with replacements
	return &AbsPaths{
		Storage:    replacePath(paths.Storage),
		Collection: replacePath(paths.Collection),
		Entities:   replacePath(paths.Entities),
		Cache:      replacePath(paths.Cache),
	}, nil
}
