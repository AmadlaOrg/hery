package storage

import (
	"fmt"
	"github.com/AmadlaOrg/hery/util/file"
	"os"
	"path/filepath"
	"runtime"
)

type Storage interface {
	Paths(collectionName string) (*AbsPaths, error)
	Main() (string, error)
	EntityPath(collectionPath, entityRelativePath string) string
}

type AbsPaths struct {
	Storage    string
	Collection string
	Entities   string
	Cache      string
}

var (
	osGetenv     = os.Getenv
	osGetwd      = os.Getwd
	filepathAbs  = filepath.Abs
	filepathJoin = filepath.Join
	fileExists   = file.Exists
)

// Paths returns the absolute path to where the entities are stored
func (d *AbsPaths) Paths(collectionName string) (*AbsPaths, error) {
	mainPath, err := d.Main()
	if err != nil {
		return &AbsPaths{}, err
	}

	collectionPath := d.collectionPath(mainPath, collectionName)
	entityPath := d.entitiesPath(collectionPath)
	cachePath := d.cachePath(collectionName, collectionPath)

	return &AbsPaths{
		Storage:    mainPath,
		Collection: collectionPath,
		Entities:   entityPath,
		Cache:      cachePath,
	}, nil
}

// Main returns the main path for .hery storage path
func (d *AbsPaths) Main() (string, error) {
	//
	// Using env var
	//

	envStoragePathValue := osGetenv(HeryStoragePath)

	if envStoragePathValue != "" {
		envStoragePath, err := filepathAbs(envStoragePathValue)
		if err != nil {
			return "", err
		}
		return envStoragePath, nil
	}

	//
	// Using current location
	//

	cwd, err := osGetwd()
	if err != nil {
		return "", err
	}

	localStoragePath := filepathJoin(cwd, ".hery")

	if fileExists(localStoragePath) {
		return localStoragePath, nil
	}

	//
	// Default
	//

	var mainDir string
	switch runtime.GOOS {
	case "windows":
		appDataDir := osGetenv("APPDATA")
		mainDir = filepathJoin(appDataDir, "Hery")
	default: // "linux" and "darwin" (macOS)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("error getting home directory: %s", err)
		}
		mainDir = filepathJoin(homeDir, ".hery")
	}

	return mainDir, nil
}

// collectionPath returns the collection absolute path
func (d *AbsPaths) collectionPath(mainPath, collectionName string) string {
	return filepathJoin(mainPath, collectionName)
}

// entityPath returns the entity absolute path
func (d *AbsPaths) entitiesPath(collectionPath string) string {
	return filepathJoin(collectionPath, "entity")
}

// cachePath returns the collection cache absolute path
func (d *AbsPaths) cachePath(collectionName, collectionPath string) string {
	return filepathJoin(collectionPath, fmt.Sprintf("%s.cache", collectionName))
}

// EntityPath returns the absolute path to a specific entity
func (d *AbsPaths) EntityPath(entitiesPath, entityRelativePath string) string {
	return filepathJoin(entitiesPath, entityRelativePath)
}
