package build

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/AmadlaOrg/hery/util/git"
	"github.com/AmadlaOrg/hery/util/url"
	"strings"
)

type Interface interface {
	MetaFromRemote(entityUri string) entity.Entity
}

type Build struct {
	Git                     git.Interface
	EntityValidation        validation.Interface
	EntityVersion           version.Interface
	EntityVersionValidation versionValidationPkg.Interface
	Storage                 storage.Storage
}

// MetaFromRemote gathers as many details about an Entity as possible from git and from the URI pass to populate the
// Entity struct
//
// It also validates values that are pass to it.
func (b *Build) MetaFromRemote(collectionName, entityUri string) (entity.Entity, error) {
	if !b.EntityValidation.EntityUrl(entityUri) {
		return entity.Entity{
			Exist: false,
		}, errors.New("invalid entity url")
	}

	paths, err := b.Storage.Paths(collectionName)
	if err != nil {
		return entity.Entity{}, err
	}

	var entityUrlPath string
	var entityFullRepoUrl string
	var entityVersion string
	//var uriEntityVersion string
	if strings.Contains(entityUri, "@") {
		uriEntityVersion, err := b.EntityVersion.Extract(entityUri)
		if err != nil {
			return entity.Entity{
				Exist: false,
			}, fmt.Errorf("error extracting version: %v", err)
		}

		entityUrlPath = url.EntityPathUrl(entityUri, uriEntityVersion)
		entityFullRepoUrl = url.EntityFullRepoUrl(entityUrlPath)

		entityVersionList, err := b.EntityVersion.List(entityFullRepoUrl)
		if err != nil {
			return entity.Entity{
				Exist: false,
			}, fmt.Errorf("error listing versions: %v", err)
		}

		//var versionExists = false
		if uriEntityVersion == "latest" {

			entityVersion, err = b.EntityVersion.Latest(entityVersionList)
			if err != nil {
				return entity.Entity{
					Exist: false,
				}, fmt.Errorf("error finding latest version: %v", err)
			}
		} else if !b.EntityVersionValidation.Format(uriEntityVersion) {
			return entity.Entity{
				Exist: false,
			}, nil
		}
		/*if b.EntityVersionValidation.Exists(entityVersion, entityVersionList) {

		}*/
		// TODO: Check with git if the version actually exists
	}

	entityPath := paths.EntityPath(paths.Collection, entityUrlPath)

	return entity.Entity{
		Name:    "",
		Uri:     entityUri,
		Origin:  entityUrlPath,
		Version: entityVersion,
		AbsPath: entityPath,
		Have:    false,
		Hash:    "",
		Exist:   true,
	}, nil
}

// MetaFromLocal gathers as many details about an Entity as possible from the local storage and from the URI pass to
// populate the Entity struct
//
// It also validates values that are pass to it and what is set in storage.
func (b *Build) MetaFromLocal(entityUri string) entity.Entity {

	return entity.Entity{
		Name:    "",
		Uri:     entityUri,
		Origin:  "",
		Version: "",
		AbsPath: "",
		Have:    false,
		Hash:    "",
		Exist:   true,
	}
}
