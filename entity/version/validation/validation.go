package validation

import (
	"errors"
	versionPkg "github.com/AmadlaOrg/hery/entity/version"
	"regexp"
)

// Exists checks if a specific version exists in the list of versions
func Exists(entityUrlPath, version string) (bool, error) {
	versions, err := versionPkg.List(entityUrlPath)
	if err != nil {
		return false, err
	}

	for _, v := range versions {
		if v == version {
			return true, nil
		}
	}

	return false, nil
}

// Format validates that the version follows one of these formats: `v1.0.0` or `v1.0` or `v1`
func Format(version string) bool {
	re := regexp.MustCompile(versionPkg.Format)

	if re.MatchString(version) {
		return true
	}

	return false
}

// PseudoFormat
func PseudoFormat(pseudoVersion string) (string, error) {
	re := regexp.MustCompile(pseudoVersion)
	if re.MatchString(pseudoVersion) {
		return pseudoVersion, nil
	}

	return "", errors.New("invalid pseudo-version formatting")
}
