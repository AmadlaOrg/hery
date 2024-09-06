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

// IBuild to help with mocking and to gather metadata from remote and local sources.
type IBuild interface {
	Meta(paths storage.AbsPaths, entityUri string) (entity.Entity, error)
	metaFromLocalWithVersion(entityUri, entityVersion string) (entity.Entity, error)
	metaFromRemoteWithoutVersion(entityUri string) (entity.Entity, error)
	metaFromRemoteWithVersion(entityUri, entityVersion string) (entity.Entity, error)
	constructOrigin(entityUri, name, version string) string
}

// SBuild struct implements the MetaBuilder interface.
type SBuild struct {
	Git                     git.IGit
	Entity                  entity.IEntity
	EntityValidation        validation.IValidation
	EntityVersion           version.IVersion
	EntityVersionValidation versionValidationPkg.IValidation
}

// Help with mocking
var (
	uuidNew = uuid.New
)

// Meta gathers as many details about an Entity as possible from git and from the URI passed to populate the
// Entity struct. It also validates values that are passed to it.
func (s *SBuild) Meta(paths storage.AbsPaths, entityUri string) (entity.Entity, error) {
	var (
		entityVals = entity.Entity{
			Have:  false,
			Exist: false,
		}
		err error
	)

	if !s.EntityValidation.EntityUri(entityUri) {
		return entityVals, errors.New("invalid entity url")
	}

	dir, err := s.Entity.FindEntityDir(paths, entityVals)
	if !errors.Is(err, entity.ErrorNotFound) &&
		!errors.Is(err, entity.ErrorMultipleFound) &&
		err != nil {
		return entityVals, err
	} else if err == nil {
		entityVals.AbsPath = dir
		entityVals.Have = true
	}

	if strings.Contains(entityUri, "@") {
		entityVersion, err := s.EntityVersion.Extract(entityUri)
		if err != nil {
			if !errors.Is(err, version.ErrorExtractNoVersionFound) {
				return entityVals, fmt.Errorf("error extracting version: %v", err)
			} else {
				entityVersion = "latest"
			}
		}

		if errors.Is(err, entity.ErrorNotFound) || entityVersion == "latest" {
			entityVals, err = s.metaFromRemoteWithVersion(entityUri, entityVersion)
			if err != nil {
				return entityVals, err
			}
		} else {
			entityVals, err = s.metaFromLocalWithVersion(entityUri, entityVersion)
			if err != nil {
				return entityVals, err
			}
		}
	} else {
		entityVals, err = s.metaFromRemoteWithoutVersion(entityUri)
		if err != nil {
			return entityVals, err
		}
	}

	entityVals.AbsPath = filepath.Join(paths.Entities, entityVals.Entity)
	entityVals.Id = uuidNew().String()
	entityVals.Exist = true

	return entityVals, nil
}

// metaFromLocalWithVersion
func (s *SBuild) metaFromLocalWithVersion(entityUri, entityVersion string) (entity.Entity, error) {
	var (
		entityVals = entity.Entity{
			Have:  false,
			Exist: false,
		}
		err error
	)

	entityUriWithoutVersion := url.TrimVersion(entityUri, entityVersion)
	entityVals.RepoUrl, err = url.ExtractRepoUrl(entityUriWithoutVersion)
	if err != nil {
		return entityVals, fmt.Errorf("error extracting repo url: %v", err)
	}

	// TODO: Get hash
	// entityVals.Hash

	entityVals.Have = true
	entityVals.Exist = true
	entityVals.IsPseudoVersion = false
	entityVals.Name = filepath.Base(entityUriWithoutVersion)
	entityVals.Version = entityVersion
	entityVals.Entity = entityUri
	entityVals.Origin = s.constructOrigin(entityVals.Entity, entityVals.Name, entityVals.Version)

	return entityVals, nil
}

// metaFromRemoteWithoutVersion
func (s *SBuild) metaFromRemoteWithoutVersion(entityUri string) (entity.Entity, error) {
	var (
		entityVals = entity.Entity{
			Have:  false,
			Exist: false,
		}
		entityVersion string
	)

	repoUrl, err := url.ExtractRepoUrl(entityUri)
	if err != nil {
		return entityVals, fmt.Errorf("error extracting repo url: %v", err)
	}
	entityVals.RepoUrl = repoUrl

	entityVersionList, err := s.EntityVersion.List(entityVals.RepoUrl)
	if err != nil {
		return entityVals, fmt.Errorf("error listing versions: %v", err)
	}

	if len(entityVersionList) == 0 {
		entityVersion, err = s.EntityVersion.GeneratePseudo(entityVals.RepoUrl)
		if err != nil {
			return entityVals, err
		}
		entityVals.IsPseudoVersion = true
	} else {
		entityVersion, err = s.EntityVersion.Latest(entityVersionList)
		if err != nil {
			return entityVals, err
		}
		entityVals.IsPseudoVersion = false
	}

	entityVals.Name = filepath.Base(entityUri)
	entityVals.Version = entityVersion
	entityVals.Entity = fmt.Sprintf("%s@%s", entityUri, entityVersion)
	entityVals.Origin = s.constructOrigin(entityVals.Entity, entityVals.Name, entityVals.Version)

	return entityVals, nil
}

// metaFromRemoteWithVersion
func (s *SBuild) metaFromRemoteWithVersion(entityUri, entityVersion string) (entity.Entity, error) {
	var (
		entityVals = entity.Entity{
			Have:  false,
			Exist: false,
		}
		err error
	)

	entityUriWithoutVersion := url.TrimVersion(entityUri, entityVersion)
	entityVals.RepoUrl, err = url.ExtractRepoUrl(entityUriWithoutVersion)
	if err != nil {
		return entityVals, fmt.Errorf("error extracting repo url: %v", err)
	}

	entityVersionList, err := s.EntityVersion.List(entityVals.RepoUrl)
	if err != nil {
		return entityVals, fmt.Errorf("error listing versions: %v", err)
	}

	if entityVersion == "latest" {
		if len(entityVersionList) == 0 {
			entityVersion, err = s.EntityVersion.GeneratePseudo(entityVals.RepoUrl)
			if err != nil {
				return entityVals, err
			}
			entityVals.IsPseudoVersion = true
		} else {
			entityVersion, err = s.EntityVersion.Latest(entityVersionList)
			if err != nil {
				return entityVals, fmt.Errorf("error finding latest version: %v", err)
			}
			entityVals.IsPseudoVersion = false
		}
	} else if !s.EntityVersionValidation.Format(entityVersion) &&
		!s.EntityVersionValidation.PseudoFormat(entityVersion) {
		return entityVals, fmt.Errorf("invalid entity version format: %v", entityVersion)
	} else if !s.EntityVersionValidation.Exists(entityVersion, entityVersionList) {
		return entityVals, fmt.Errorf("invalid entity version: %v", entityVersion)
	}

	entityVals.Name = filepath.Base(entityUriWithoutVersion)
	entityVals.Version = entityVersion
	entityVals.Entity = entityUri
	entityVals.Origin = s.constructOrigin(entityVals.Entity, entityVals.Name, entityVals.Version)

	return entityVals, nil
}

// constructOrigin
func (s *SBuild) constructOrigin(entityUri, name, version string) string {
	return strings.Replace(
		entityUri,
		fmt.Sprintf("%s@%s", name, version),
		"",
		1)
}
