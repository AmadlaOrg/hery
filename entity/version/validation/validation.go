package validation

import (
	"github.com/AmadlaOrg/LibraryUtils/str"
	versionPkg "github.com/AmadlaOrg/hery/entity/version"
	"regexp"
)

// IValidation is an interface for version validation.
type IValidation interface {
	Exists(version string, versions []string) bool
	Format(version string) bool
	PseudoFormat(pseudoVersion string) bool
}

// SValidation struct implements the VersionValidator interface.
type SValidation struct {
	Version versionPkg.IVersion
}

// Exists checks if a specific version exists in the list of versions.
func (s *SValidation) Exists(version string, versions []string) bool {
	return str.ExistInStringArr(version, versions)
}

// Format validates that the version follows one of these formats: `v1.0.0`, `v1.0`, or `v1`.
func (s *SValidation) Format(version string) bool {
	return regexp.MustCompile(versionPkg.Format).MatchString(version)
}

// PseudoFormat validates that the pseudo version follows a specified format.
func (s *SValidation) PseudoFormat(pseudoVersion string) bool {
	// Define a regex pattern for pseudo versions.
	pseudoVersionPattern := versionPkg.PseudoVersionFormat
	re := regexp.MustCompile(pseudoVersionPattern)
	return re.MatchString(pseudoVersion)
}

// AnyFormat validates that the version pass is in standard format or pseudo format or `latest` word.
func (s *SValidation) AnyFormat(version string) bool {
	if !s.Format(version) && !s.PseudoFormat(version) && version != "latest" {
		return false
	}
	return true
}
