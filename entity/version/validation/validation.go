package validation

import (
	versionPkg "github.com/AmadlaOrg/hery/entity/version"
	"github.com/AmadlaOrg/hery/util/str"
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
	Version versionPkg.SVersion
}

// Exists checks if a specific version exists in the list of versions.
func (vs *VersionValidation) Exists(version string, versions []string) bool {
	return str.ExistInStringArr(version, versions)
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

// AnyFormat validates that the version pass is in standard format or pseudo format or `latest` word.
func (vs *VersionValidation) AnyFormat(version string) bool {
	if !vs.Format(version) && !vs.PseudoFormat(version) && version != "latest" {
		return false
	}
	return true
}
