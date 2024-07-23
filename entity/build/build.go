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
	MetaFromLocal(entityUri string) entity.Entity
}

// Builder struct implements the MetaBuilder interface.
type Builder struct {
	Git                     git.RepoManager
	EntityValidation        validation.Interface
	EntityVersion           version.Manager
	EntityVersionValidation versionValidationPkg.VersionValidation
	Storage                 storage.AbsPaths
}

// MetaFromRemote gathers as many details about an Entity as possible from git and from the URI passed to populate the
// Entity struct. It also validates values that are passed to it.
func (b *Builder) MetaFromRemote(collectionName, entityUri string) (entity.Entity, error) {
	var entityVals = entity.Entity{
		Have:  false,
		Exist: false,
	}

	if !b.EntityValidation.EntityUrl(entityUri) {
		return entityVals, errors.New("invalid entity url")
	}

	paths, err := b.Storage.Paths(collectionName)
	if err != nil {
		return entityVals, err
	}

	var entityUriWithVersion string
	var entityUrlPath string
	var entityFullRepoUrl string
	var entityVersion string
	var originPath string
	//var uriEntityVersion string
	if strings.Contains(entityUri, "@") {
		entityUriWithVersion = entityUri
		uriEntityVersion, err := b.EntityVersion.Extract(entityUri)
		if err != nil {
			return entityVals, fmt.Errorf("error extracting version: %v", err)
		}

		entityUrlPath = url.EntityPathUrl(entityUri, uriEntityVersion)
		entityFullRepoUrl = url.EntityFullRepoUrl(entityUrlPath)

		entityVersionList, err := b.EntityVersion.List(entityFullRepoUrl)
		if err != nil {
			return entityVals, fmt.Errorf("error listing versions: %v", err)
		}

		//var versionExists = false
		if uriEntityVersion == "latest" {
			entityVersion, err = b.EntityVersion.Latest(entityVersionList)
			if err != nil {
				return entityVals, fmt.Errorf("error finding latest version: %v", err)
			}
		} else if !b.EntityVersionValidation.Format(uriEntityVersion) {
			return entityVals, nil
		}

		entityVals.Uri = entityUrlPath
		entityVals.Version = uriEntityVersion
		/*if b.EntityVersionValidation.Exists(entityVersion, entityVersionList) {

		}*/
		// TODO: Check with git if the version actually exists
	} else {
		entityVals.Uri = url.EntityFullRepoUrl(entityUri)
		entityVersionList, err := b.EntityVersion.List(entityVals.Uri)
		if err != nil {
			return entityVals, fmt.Errorf("error listing versions: %v", err)
		}

		if len(entityVersionList) == 0 {
			entityVersion, err = b.EntityVersion.GeneratePseudo(entityFullRepoUrl)
			if err != nil {
				return entityVals, err
			}
		} else {
			entityVersion, err = b.EntityVersion.Latest(entityVersionList)
			if err != nil {
				return entityVals, err
			}
		}

		entityVals.Version = entityVersion
		entityUrlPath = entityUri
		entityUriWithVersion = fmt.Sprintf("%s@%s", entityFullRepoUrl, entityVersion)

	}

	entityPath := paths.EntityPath(paths.Entities, entityUrlPath)
	entityDirName := filepath.Base(entityPath)

	return entity.Entity{
		Id:      uuid.New().String(),
		Entity:  entityUriWithVersion,
		Name:    "",
		DirName: entityDirName,
		Uri:     entityUri,
		Origin:  originPath,
		Version: entityVersion,
		AbsPath: entityPath,
		Have:    false,
		Hash:    "",
		Exist:   true,
	}, nil
}

// MetaFromLocal gathers as many details about an Entity as possible from the local storage and from the URI passed to
// populate the Entity struct. It also validates values that are passed to it and what is set in storage.
func (b *Builder) MetaFromLocal(entityUri string) entity.Entity {
	// Implementation logic for MetaFromLocal

	return entity.Entity{
		Name:    "",
		Uri:     entityUri,
		Id:      uuid.New().String(),
		Origin:  "",
		Version: "",
		AbsPath: "",
		Have:    false,
		Hash:    "",
		Exist:   true,
	}
}
