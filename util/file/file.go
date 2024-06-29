package file

import "os"

type File interface {
	Exists(path string) bool
}

// Exists verify that a file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
