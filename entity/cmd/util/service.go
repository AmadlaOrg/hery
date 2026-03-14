package util

import "github.com/AmadlaOrg/hery/storage"

// New to set up the Util service
func New() Util {
	return &utilImpl{
		New: storage.New(),
	}
}
