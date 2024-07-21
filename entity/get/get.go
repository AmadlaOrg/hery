package get

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/AmadlaOrg/hery/util/git"
)

// EntityGetter is an interface for getting entities.
type EntityGetter interface {
	Get(collectionName, storagePath string, args []string) error
	download(entityUrls []string, collectionStoragePath string) error
}

// GetterService struct implements the EntityGetter interface.
type GetterService struct {
	Git                     git.RepoManager
	EntityValidation        validation.Interface
	EntityVersion           *version.Service
	EntityVersionValidation versionValidationPkg.VersionValidator
}

// Get retrieves entities based on the provided collection name and arguments.
func (gs *GetterService) Get(collectionName string, storagePath *storage.AbsPaths, args []string) error {
	// Validate that all the URLs passed in the arguments are valid
	if len(args) == 0 {
		return errors.New("no entity URL(s) specified")
	}

	for _, arg := range args {
		if !gs.EntityValidation.EntityUrl(arg) {
			return fmt.Errorf("invalid entity URL: %s", arg)
		}
	}

	return nil
	//return gs.download(args, storagePath)
}

// download retrieves entities in parallel.
/*func (gs *GetterService) download(entityUrls []string, collectionStoragePath string) error {
var wg sync.WaitGroup
wg.Add(len(entityUrls))

errCh := make(chan error, len(entityUrls)) // Channel to collect errors
//TODO: entityPaths := make([]string, len(entityUrls)) Maybe check with directories

for _, entityUrl := range entityUrls {
	go func(entityUrl string) {
		defer wg.Done()

		// TODO: skip (continue) loop-iteration if entity with same version was already downloaded/installed maybe make a Map

		//var entityUrlPath string
		var entityFullRepoUrl string
		var entityVersion string

		if strings.Contains(entityUrl, "@") {
			/*entityUrlPath, entityFullRepoUrl, entityVersion, err := gs.processEntityUrlWithVersion(entityUrl)
			if err != nil {
				errCh <- err
				return
			}*/
/*} else {
/*entityUrlPath, entityFullRepoUrl, entityVersion, err := gs.processEntityUrlWithoutVersion(entityUrl)
if err != nil {
	errCh <- err
	return
}*/
/*	entityUrl = fmt.Sprintf("%s@%s", entityUrl, entityVersion)
			}

			destination := filepath.Join(collectionStoragePath, entityUrl)

			// Download the Entity with `git clone`
			if err := gs.Git.FetchRepo(entityFullRepoUrl, destination); err != nil {
				errCh <- fmt.Errorf("error fetching repo: %v", err)
			}

			println(entityVersion)
		}(entityUrl)
	}

	wg.Wait()
	close(errCh)

	var combinedErr error
	for e := range errCh {
		if combinedErr == nil {
			combinedErr = e
		} else {
			combinedErr = fmt.Errorf("%v; %v", combinedErr, e)
		}
	}

	return combinedErr
}

func (gs *GetterService) processEntityUrlWithVersion(entityUrl string) (string, string, string, error) {
	uriEntityVersion, err := gs.EntityVersion.Extract(entityUrl)
	if err != nil {
		return "", "", "", fmt.Errorf("error extracting version: %v", err)
	}

	entityUrlPath := url.EntityPathUrl(entityUrl, uriEntityVersion)
	entityFullRepoUrl := url.EntityFullRepoUrl(entityUrlPath)

	entityVersionList, err := gs.EntityVersion.List(entityFullRepoUrl)
	if err != nil {
		return "", "", "", fmt.Errorf("error listing versions: %v", err)
	}

	var entityVersion string
	var versionExists bool
	if uriEntityVersion == "latest" {
		entityVersion, err = gs.EntityVersion.Latest(entityVersionList)
		if err != nil {
			return "", "", "", fmt.Errorf("error finding latest version: %v", err)
		}
		versionExists = true
	} else if !gs.EntityVersionValidation.Format(uriEntityVersion) {
		return "", "", "", errors.New("entity version in the entity url is wrong format")
	}

	//if !versionExists {
	// Check if the version exists
	/*versionExists := gs.EntityVersionValidation.Exists(entityUrlPath, uriEntityVersion)
	if !versionExists {
		return "", "", "", fmt.Errorf("the version of the entity URL does not exist: %s", entityUrl)
	}*/
//}
/*
	fmt.Println(versionExists)

	return entityUrlPath, entityFullRepoUrl, entityVersion, nil
}

func (gs *GetterService) processEntityUrlWithoutVersion(entityUrl string) (string, string, string, error) {
	entityUrlPath := entityUrl
	entityFullRepoUrl := url.EntityFullRepoUrl(entityUrlPath)
	entityVersionList, err := gs.EntityVersion.List(entityFullRepoUrl)
	if err != nil {
		return "", "", "", fmt.Errorf("error listing versions: %v", err)
	}

	var entityVersion string
	if len(entityVersionList) == 0 {
		entityVersion, err = gs.EntityVersion.GeneratePseudo(entityFullRepoUrl)
		if err != nil {
			return "", "", "", fmt.Errorf("error generating pseudo version: %v", err)
		}
	} else {
		entityVersion, err = gs.EntityVersion.Latest(entityVersionList)
		if err != nil {
			return "", "", "", fmt.Errorf("error finding latest version: %v", err)
		}
	}

	return entityUrlPath, entityFullRepoUrl, entityVersion, nil
}*/
