package entity

import (
	"errors"
	"fmt"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/errtypes"
	"github.com/AmadlaOrg/hery/storage"
	"os"
	"path/filepath"
)

// FindEntityDir can find pseudo versioned entity directories and static versioned entities.
func FindEntityDir(paths storage.AbsPaths, entityVals Entity) (string, error) {
	entityVersionValidation := versionValidationPkg.NewEntityVersionValidationService()
	if !entityVersionValidation.PseudoFormat(entityVals.Version) {
		exactPath := entityVals.Entity

		// Check if the directory exists
		if _, err := os.Stat(exactPath); os.IsNotExist(err) {
			return "", errors.Join(
				errtypes.NotFoundError,
				fmt.Errorf("no matching directory found for exact version: %s", exactPath))
		} else if err != nil {
			return "", err
		}

		// Return the exact path if it exists
		return exactPath, nil
	}

	// Construct the pattern
	pattern := filepath.Join(
		paths.Entities, entityVals.Origin, GeneratePseudoVersionPattern(entityVals.Name, entityVals.Version))

	// Use Glob to find directories matching the pattern
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", errors.Join(
			errtypes.NotFoundError,
			fmt.Errorf("no matching directories found for pattern: %s", pattern))
	}

	if len(matches) > 1 {
		return "", errors.Join(
			errtypes.MultipleFoundError,
			fmt.Errorf("multiple matching directories found for pattern: %s", pattern))
	}

	// Return the matched directory
	return matches[0], nil
}

// CheckDuplicateEntity checks if entityMeta is already in entityBuilds.
func CheckDuplicateEntity(entities []Entity, entityMeta Entity) error {
	entityVersionValidation := versionValidationPkg.NewEntityVersionValidationService()

	for _, existingEntity := range entities {
		if existingEntity.Origin == entityMeta.Origin &&
			existingEntity.Name == entityMeta.Name {
			if entityVersionValidation.PseudoFormat(existingEntity.Version) &&
				entityVersionValidation.PseudoFormat(entityMeta.Version) {
				// Check pseudo versions
				if GeneratePseudoVersionPattern(existingEntity.Name, existingEntity.Version) == GeneratePseudoVersionPattern(entityMeta.Name, entityMeta.Version) {
					return fmt.Errorf("duplicate entity found: %v", entityMeta)
				}
			} else if existingEntity.Version == entityMeta.Version {
				// Exact version match
				return fmt.Errorf("duplicate entity found: %v", entityMeta)
			}
		}
	}
	return nil
}

// GeneratePseudoVersionPattern generates a pattern string for pseudo-versioned entities based on their name and version.
func GeneratePseudoVersionPattern(name, version string) string {
	return fmt.Sprintf("%s@%s-*-%s", name, version[:6], version[22:])
}
