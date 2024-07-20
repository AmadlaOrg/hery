package validation

import (
	versionPkg "github.com/AmadlaOrg/hery/entity/version"
	utilString "github.com/AmadlaOrg/hery/util/string"
	"regexp"
)

// VersionValidator is an interface for version validation.
type VersionValidator interface {
	Exists(version string, versions []string) bool
	Format(version string) bool
	PseudoFormat(pseudoVersion string) bool
}

// VersionValidation struct implements the VersionValidator interface.
type VersionValidation struct {
	Version versionPkg.Service
}

// Exists checks if a specific version exists in the list of versions.
func (vs *VersionValidation) Exists(version string, versions []string) bool {
	return utilString.ExistInStringArr(version, versions)
}

// Format validates that the version follows one of these formats: `v1.0.0`, `v1.0`, or `v1`.
func (vs *VersionValidation) Format(version string) bool {
	return regexp.MustCompile(versionPkg.Format).MatchString(version)
}

// PseudoFormat validates that the pseudo version follows a specified format.
func (vs *VersionValidation) PseudoFormat(pseudoVersion string) bool {
	// Define a regex pattern for pseudo versions.
	pseudoVersionPattern := versionPkg.PseudoVersionFormat
	re := regexp.MustCompile(pseudoVersionPattern)
	return re.MatchString(pseudoVersion)
}
