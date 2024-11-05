package util

// NewEntityCmdUtilService to set up the Util service
func NewEntityCmdUtilService() IUtil {
	return &SUtil{}
}
