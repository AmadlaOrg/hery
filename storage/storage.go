package storage

import (
	"fmt"
	"github.com/AmadlaOrg/LibraryUtils/file"
	"os"
	"path/filepath"
	"runtime"
)

type Storage interface {
	Paths() (*AbsPaths, error)
	Main() (string, error)
	EntityPath(entitiesPath, entityRelativePath string) string
	TmpPaths() (*AbsPaths, error)
	TmpMain() (string, error)
	MakePaths(paths AbsPaths) error
}

// AbsPaths holds the absolute paths for hery storage.
// Draft 3.2 layout: no collection directories.
type AbsPaths struct {
	Storage  string // Global cache root: ~/.cache/hery/
	Entities string // Entity type cache: ~/.cache/hery/entity/
	Cache    string // Project cache: .hery.cache (in project root)
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

// Paths returns the absolute paths for hery storage
func (d *AbsPaths) Paths() (*AbsPaths, error) {
	mainPath, err := d.Main()
	if err != nil {
		return &AbsPaths{}, err
	}

	entitiesPath := filepathJoin(mainPath, "entity")

	return &AbsPaths{
		Storage:  mainPath,
		Entities: entitiesPath,
	}, nil
}

// Main returns the root path for hery's global cache
func (d *AbsPaths) Main() (string, error) {
	// Using env var
	envStoragePathValue := os.Getenv(HeryStoragePath)
	if envStoragePathValue != "" {
		envStoragePath, err := filepathAbs(envStoragePathValue)
		if err != nil {
			return "", err
		}
		return envStoragePath, nil
	}

	// Default: ~/.cache/hery/
	var mainDir string
	switch runtime.GOOS {
	case "windows":
		appDataDir := os.Getenv("LOCALAPPDATA")
		if appDataDir == "" {
			appDataDir = os.Getenv("APPDATA")
		}
		mainDir = filepathJoin(appDataDir, "hery")
	default: // "linux" and "darwin" (macOS)
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("error getting home directory: %s", err)
			}
			cacheDir = filepathJoin(homeDir, ".cache")
		}
		mainDir = filepathJoin(cacheDir, "hery")
	}

	return mainDir, nil
}

// EntityPath returns the absolute path to a specific entity
func (d *AbsPaths) EntityPath(entitiesPath, entityRelativePath string) string {
	return filepathJoin(entitiesPath, entityRelativePath)
}

// TmpPaths returns temporary absolute paths for storage
func (d *AbsPaths) TmpPaths() (*AbsPaths, error) {
	mainTmpPath, err := d.TmpMain()
	if err != nil {
		return &AbsPaths{}, err
	}

	entitiesPath := filepathJoin(mainTmpPath, "entity")

	return &AbsPaths{
		Storage:  mainTmpPath,
		Entities: entitiesPath,
	}, nil
}

// TmpMain returns a temporary main path for hery storage
func (d *AbsPaths) TmpMain() (string, error) {
	tempDir, err := osMkdirTemp("", "hery_*")
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

// MakePaths creates all storage subdirectories
func (d *AbsPaths) MakePaths(paths AbsPaths) error {
	if err := osMkdirAll(paths.Storage, perm); err != nil {
		return err
	}

	if err := osMkdirAll(paths.Entities, perm); err != nil {
		return err
	}

	return nil
}
