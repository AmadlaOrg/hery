package validation

import (
	"github.com/AmadlaOrg/LibraryUtils/str"
	versionPkg "github.com/AmadlaOrg/hery/entity/version"
	"regexp"
)

// Validator is an interface for version validation.
type Validator interface {
	Exists(version string, versions []string) bool
	Format(version string) bool
	PseudoFormat(pseudoVersion string) bool
}

// validator struct implements the VersionValidator interface.
type validator struct {
	Version versionPkg.Version
}

// Exists checks if a specific version exists in the list of versions.
func (s *validator) Exists(version string, versions []string) bool {
	return str.ExistInStringArr(version, versions)
}

// Format validates that the version follows one of these formats: `v1.0.0`, `v1.0`, or `v1`.
func (s *validator) Format(version string) bool {
	return regexp.MustCompile(versionPkg.Format).MatchString(version)
}

// PseudoFormat validates that the pseudo version follows a specified format.
func (s *validator) PseudoFormat(pseudoVersion string) bool {
	// Define a regex pattern for pseudo versions.
	pseudoVersionPattern := versionPkg.PseudoVersionFormat
	re := regexp.MustCompile(pseudoVersionPattern)
	return re.MatchString(pseudoVersion)
}

// AnyFormat validates that the version pass is in standard format or pseudo format or `latest` word.
func (s *validator) AnyFormat(version string) bool {
	if !s.Format(version) && !s.PseudoFormat(version) && version != "latest" {
		return false
	}
	return true
}
