package file

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

var (
	osStat       = os.Stat
	osIsNotExist = os.IsNotExist
	osOpen       = func(name string) (IFile, error) {
		return os.Open(name)
	}
	bytesEqual = bytes.Equal
)

// Exists verifies that a file or directory exists
func Exists(path string) bool {
	_, err := osStat(path)
	return err == nil || !osIsNotExist(err)
}

// IsFile verifies that the path points to a file and not a directory
func IsFile(path string) (bool, error) {
	// Check if the file exists and is a regular file
	info, err := osStat(path)
	if err != nil {
		// Return error if the file doesn't exist or there is an issue with the path
		return false, errors.Join(ErrorFailedToStatFile, fmt.Errorf("path %s: %v", path, err))
	}

	if info.IsDir() {
		// Return an error if the path points to a directory
		return false, errors.Join(ErrorNotAFile, fmt.Errorf("the path %s is a directory", path))
	}

	return true, nil
}

// IsValidMagic validates that the magic head matches what is in a file
func IsValidMagic(path string, magic []byte) (bool, error) {
	if ok, err := IsFile(path); !ok {
		return false, err
	}

	// Check if the file is a valid SQLite3 file by reading the first 4 bytes
	file, err := osOpen(path)
	if err != nil {
		return false, err
	}
	defer func(file IFile) {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close file %s: %v", path, closeErr)
		}
	}(file)

	// Read the first 4 bytes of the file
	header := make([]byte, 4)
	_, err = file.Read(header)
	if err != nil {
		return false, err
	}

	// Check if the file starts with the SQLite magic header
	if !bytesEqual(header, magic) {
		return false, errors.Join(ErrorMagicIsNotAMatch, fmt.Errorf("file %s does not match magic header", path))
	}

	return true, nil
}
