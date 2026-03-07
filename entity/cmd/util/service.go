package util

import "github.com/AmadlaOrg/hery/storage"

// NewEntityCmdUtilService to set up the Util service
func NewEntityCmdUtilService() IUtil {
	return &SUtil{
		NewStorageService: storage.NewStorageService(),
	}
}
