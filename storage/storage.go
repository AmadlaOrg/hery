package storage

import (
	"fmt"
	"github.com/AmadlaOrg/hery/util/file"
	"os"
	"path/filepath"
	"runtime"
)

type IStorage interface {
	Paths(collectionName string) (*AbsPaths, error)
	Main() (string, error)
	EntityPath(collectionPath, entityRelativePath string) string
	TmpPaths(collectionName string) (*AbsPaths, error)
	TmpMain() (string, error)
	MakePaths(paths AbsPaths) error
}

type AbsPaths struct {
	Storage    string // e.g.: /home/user/.hery/
	Catalog    string // e.g.: /home/user/.hery/collection/
	Collection string // e.g.: /home/user/.hery/collection/amadla/
	Entities   string // e.g.: /home/user/.hery/collection/amadla/entity/
	Cache      string // e.g.: /home/user/.hery/collection/amadla/amadla.cache
}

const perm os.FileMode = os.ModePerm

var (
	osGetwd      = os.Getwd
	filepathAbs  = filepath.Abs
	filepathJoin = filepath.Join
	fileExists   = file.Exists
	osMkdirAll   = os.MkdirAll
	osMkdirTemp  = os.MkdirTemp
)

// Paths return the absolute paths for the different parts of storage
func (d *AbsPaths) Paths(collectionName string) (*AbsPaths, error) {
	mainPath, err := d.Main()
	if err != nil {
		return &AbsPaths{}, err
	}

	return d.paths(mainPath, collectionName)
}

// Main returns the main path for `.hery` storage path
// TODO: Maybe name it Root
func (d *AbsPaths) Main() (string, error) {
	//
	// Using env var
	//

	envStoragePathValue := os.Getenv(HeryStoragePath)

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
		appDataDir := os.Getenv("APPDATA")
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

// EntityPath returns the absolute path to a specific entity
func (d *AbsPaths) EntityPath(entitiesPath, entityRelativePath string) string {
	return filepathJoin(entitiesPath, entityRelativePath)
}

// TmpPaths returns the tmp absolute paths for the different parts of storage
func (d *AbsPaths) TmpPaths(collectionName string) (*AbsPaths, error) {
	mainTmpPath, err := d.TmpMain()
	if err != nil {
		return &AbsPaths{}, err
	}

	return d.paths(mainTmpPath, collectionName)
}

// TmpMain returns the tmp main path for `.hery` storage path
func (d *AbsPaths) TmpMain() (string, error) {
	tempDir, err := osMkdirTemp("", "hery_*")
	if err != nil {
		return "", err
	}

	storageTmpPath := filepath.Join(tempDir, ".hery")
	err = osMkdirAll(storageTmpPath, perm)
	if err != nil {
		return "", err
	}

	return storageTmpPath, nil
}

// MakePaths makes all the storage subdirectories
func (d *AbsPaths) MakePaths(paths AbsPaths) error {
	err := osMkdirAll(paths.Storage, perm)
	if err != nil {
		return err
	}

	err = osMkdirAll(paths.Catalog, perm)
	if err != nil {
		return err
	}

	err = osMkdirAll(paths.Collection, perm)
	if err != nil {
		return err
	}

	err = osMkdirAll(paths.Entities, perm)
	if err != nil {
		return err
	}

	return nil
}

// paths internal function to generate the paths for different part of storage based on main path and collection name
func (d *AbsPaths) paths(mainPath, collectionName string) (*AbsPaths, error) {
	catalogPath := d.catalogPath(mainPath)
	collectionPath := d.collectionPath(catalogPath, collectionName)
	entityPath := d.entitiesPath(collectionPath)
	cachePath := d.cachePath(collectionName, collectionPath)

	return &AbsPaths{
		Storage:    mainPath,
		Catalog:    catalogPath,
		Collection: collectionPath,
		Entities:   entityPath,
		Cache:      cachePath,
	}, nil
}

// catalogPath returns the catalog absolute path
func (d *AbsPaths) catalogPath(mainPath string) string {
	return filepathJoin(mainPath, "collection")
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
