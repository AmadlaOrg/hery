package build

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/util/git"
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
}

// MetaFromRemote gathers as many details about an Entity as possible from git and from the URI pass to populate the
// Entity struct
//
// It also validates values that are pass to it.
func (b *Build) MetaFromRemote(entityUri string) (entity.Entity, error) {
	if !b.EntityValidation.EntityUrl(entityUri) {
		return entity.Entity{}, errors.New("invalid entity url")
	}

	var uriEntityVersion string
	if strings.Contains(entityUri, "@") {
		uriEntityVersion, err := b.EntityVersion.Extract(entityUri)
		if err != nil {
			return entity.Entity{}, fmt.Errorf("error extracting version: %v", err)
		}
		if !b.EntityVersionValidation.Format(uriEntityVersion) {
			return entity.Entity{}, nil
		}
		// TODO: Check with git if the version actually exists
	}

	return entity.Entity{
		Name:    "",
		Uri:     entityUri,
		Origin:  "",
		Version: uriEntityVersion,
		AbsPath: "",
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
