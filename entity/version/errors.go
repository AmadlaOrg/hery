package version

import "errors"

var (
	ErrorExtractVersionAnnotationCountMoreThanOne = errors.New("invalid URI, contains more than one '@'")
	ErrorExtractNoVersionFound                    = errors.New("no version found in URI")
	ErrorExtractInvalidVersion                    = errors.New("invalid version found in URI")
	ErrorListGitRemoteTags                        = errors.New("error listing git remote tags")
	ErrorLatestVersionsLenIsZero                  = errors.New("versions is empty")
)
