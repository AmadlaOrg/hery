package build

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"path/filepath"
	"strings"

	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/AmadlaOrg/hery/util/git"
	"github.com/AmadlaOrg/hery/util/url"
)

// MetaBuilder to help with mocking and to gather metadata from remote and local sources.
type MetaBuilder interface {
	MetaFromRemote(collectionName, entityUri string) (entity.Entity, error)
}

// Builder struct implements the MetaBuilder interface.
type Builder struct {
	Git                     git.RepoManager
	EntityValidation        validation.Interface
	EntityVersion           version.Manager
	EntityVersionValidation versionValidationPkg.VersionValidation
}

// MetaFromRemote gathers as many details about an Entity as possible from git and from the URI passed to populate the
// Entity struct. It also validates values that are passed to it.
func (b *Builder) MetaFromRemote(paths storage.AbsPaths, entityUri string) (entity.Entity, error) {
	var entityVals = entity.Entity{
		Have:  false,
		Exist: false,
	}

	if !b.EntityValidation.EntityUrl(entityUri) {
		return entityVals, errors.New("invalid entity url")
	}

	var entityVersion string
	if strings.Contains(entityUri, "@") {
		entityVersion, err := b.EntityVersion.Extract(entityUri)
		if err != nil {
			return entityVals, fmt.Errorf("error extracting version: %v", err)
		}

		entityVals.RepoUrl, err = url.ExtractRepoUrl(entityUri)
		if err != nil {
			return entityVals, fmt.Errorf("error extracting repo url: %v", err)
		}

		entityVersionList, err := b.EntityVersion.List(entityVals.RepoUrl)
		if err != nil {
			return entityVals, fmt.Errorf("error listing versions: %v", err)
		}

		if entityVersion == "latest" {
			entityVersion, err = b.EntityVersion.Latest(entityVersionList)
			if err != nil {
				return entityVals, fmt.Errorf("error finding latest version: %v", err)
			}
		} else if !b.EntityVersionValidation.Format(entityVersion) {
			return entityVals, fmt.Errorf("invalid entity version: %v", entityVersion)
		} else if !b.EntityVersionValidation.Exists(entityVersion, entityVersionList) {
			return entityVals, fmt.Errorf("invalid entity version: %v", entityVersion)
		}

		entityVals.Entity = entityUri
	} else {
		repoUrl, err := url.ExtractRepoUrl(entityUri)
		if err != nil {
			return entityVals, fmt.Errorf("error extracting repo url: %v", err)
		}
		entityVals.RepoUrl = repoUrl

		entityVersionList, err := b.EntityVersion.List(entityVals.RepoUrl)
		if err != nil {
			return entityVals, fmt.Errorf("error listing versions: %v", err)
		}

		if len(entityVersionList) == 0 {
			entityVersion, err = b.EntityVersion.GeneratePseudo(entityVals.RepoUrl)
			if err != nil {
				return entityVals, err
			}
		} else {
			entityVersion, err = b.EntityVersion.Latest(entityVersionList)
			if err != nil {
				return entityVals, err
			}
		}

		entityVals.Entity = fmt.Sprintf("%s@%s", entityUri, entityVersion)
	}

	entityVals.Name = filepath.Base(entityUri)
	entityVals.Version = entityVersion
	entityVals.Origin = strings.Replace(
		entityVals.Entity,
		fmt.Sprintf("%s@%s", entityVals.Name, entityVals.Version),
		"",
		1)
	entityVals.AbsPath = filepath.Join(paths.Entities, entityVals.Entity)
	entityVals.Id = uuid.New().String()
	entityVals.Exist = true

	return entityVals, nil
}
