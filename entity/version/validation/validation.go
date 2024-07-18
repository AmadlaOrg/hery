package validation

import (
	versionPkg "github.com/AmadlaOrg/hery/entity/version"
	utilString "github.com/AmadlaOrg/hery/util/string"
	"regexp"
)

type Interface interface {
	Exists(version string, versions []string) bool
	Format(version string) bool
	PseudoFormat(pseudoVersion string) bool
}

type Validation struct {
	Version versionPkg.Interface
}

// Exists checks if a specific version exists in the list of versions
func (vs *Validation) Exists(version string, versions []string) bool {
	return utilString.ExistInStringArr(version, versions)
}

// Format validates that the version follows one of these formats: `v1.0.0` or `v1.0` or `v1`
func (vs *Validation) Format(version string) bool {
	return regexp.MustCompile(versionPkg.Format).MatchString(version)
}

// PseudoFormat
func (vs *Validation) PseudoFormat(pseudoVersion string) bool {
	// TODO: Needs a regex
	re := regexp.MustCompile(pseudoVersion)
	return re.MatchString(pseudoVersion)
}
