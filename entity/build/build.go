package build

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/hery/errtypes"
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
	MetaFromRemote(paths storage.AbsPaths, entityUri string) (entity.Entity, error)
	metaFromRemoteWithoutVersion(entityUri string) (entity.Entity, error)
	metaFromRemoteWithVersion(entityUri string) (entity.Entity, error)
}

// SBuild struct implements the MetaBuilder interface.
type SBuild struct {
	Git                     git.RepoManager
	Entity                  entity.IEntity
	EntityValidation        validation.Interface
	EntityVersion           version.IVersion
	EntityVersionValidation versionValidationPkg.VersionValidation
}

// MetaFromRemote gathers as many details about an Entity as possible from git and from the URI passed to populate the
// Entity struct. It also validates values that are passed to it.
func (s *SBuild) MetaFromRemote(paths storage.AbsPaths, entityUri string) (entity.Entity, error) {
	var entityVals = entity.Entity{
		Have:  false,
		Exist: false,
	}

	if !s.EntityValidation.EntityUrl(entityUri) {
		return entityVals, errors.New("invalid entity url")
	}

	//var entityVersion string
	if strings.Contains(entityUri, "@") {
		entityVals, err := s.metaFromRemoteWithVersion(entityUri)
		if err != nil {
			return entityVals, err
		}
	} else {
		entityVals, err := s.metaFromRemoteWithoutVersion(entityUri)
		if err != nil {
			return entityVals, err
		}
	}

	entityVals.AbsPath = filepath.Join(paths.Entities, entityVals.Entity)
	entityVals.Id = uuid.New().String()
	entityVals.Exist = true

	dir, err := s.Entity.FindEntityDir(paths, entityVals)
	if errors.Is(err, errtypes.MultipleFoundError) {
		return entityVals, err
	} else if !errors.Is(err, errtypes.NotFoundError) && err != nil {
		return entityVals, err
	}

	if err == nil {
		entityVals.AbsPath = dir
		entityVals.Have = true
	}

	return entityVals, nil
}

func (s *SBuild) metaFromRemoteWithoutVersion(entityUri string) (entity.Entity, error) {
	var entityVals = entity.Entity{
		Have:  false,
		Exist: false,
	}
	var entityVersion string

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
	entityVals.Origin = strings.Replace(
		entityVals.Entity,
		fmt.Sprintf("%s@%s", entityVals.Name, entityVals.Version),
		"",
		1)

	return entityVals, nil
}

func (s *SBuild) metaFromRemoteWithVersion(entityUri string) (entity.Entity, error) {
	var entityVals = entity.Entity{
		Have:  false,
		Exist: false,
	}
	var entityVersion string

	entityVersion, err := s.EntityVersion.Extract(entityUri)
	if err != nil {
		if !errors.Is(err, version.ErrorExtractNoVersionFound) {
			return entityVals, fmt.Errorf("error extracting version: %v", err)
		} else {
			entityVersion = "latest"
		}
	}

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
	entityVals.Origin = strings.Replace(
		entityVals.Entity,
		fmt.Sprintf("%s@%s", entityVals.Name, entityVals.Version),
		"",
		1)

	return entityVals, nil
}
