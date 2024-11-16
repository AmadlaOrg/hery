package collection

import (
	"fmt"
	"github.com/AmadlaOrg/hery/storage"
	"os"
	"path/filepath"
)

type ICollection interface {
	Select(collectionName string) IEntityCollection
	Create(collectionName string) error
	Remove(collectionName string) error
	Exists(collectionName string) bool
	List() ([][]string, error)
}

type SCollection struct {
	Storage storage.IStorage

	// Data
	Collections *[]*Collection
}

const perm os.FileMode = os.ModePerm

var (
	osMkdirAll = os.MkdirAll
)

// Select collection to work with
func (s *SCollection) Select(collectionName string) IEntityCollection {

	return &SEntityCollection{}
}

// Create collection components (collection directory, the directories and cache file)
func (s *SCollection) Create(collectionName string) error {
	if collectionName == "" {
		return fmt.Errorf("collection name is empty")
	}
	// TODO: Add validation (regex, minimum and maximum)

	collectionPath, err := s.Storage.Main()
	if err != nil {
		return err
	}

	err = osMkdirAll(filepath.Join(collectionPath, collectionName), perm)
	if err != nil {
		return err
	}

	return nil
}

// Remove a specific collection
func (s *SCollection) Remove(collectionName string) error {
	return nil
}

// Exists
func (s *SCollection) Exists(collectionName string) bool {
	/*list, err := s.List()
	if err != nil {
		return false
	}

	/*for _, name := range list {
		if name == collectionName {
			return true
		}
	}*/
	return false
}

// List gathers all the directories that are collections in the `.henry` main storage directory to return them
func (s *SCollection) List() ([][]string, error) {
	path, err := s.Storage.Main()
	if err != nil {
		return nil, err
	}

	// Open the directory
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	// Prepare the data for the table
	var directories [][]string
	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, []string{file.Name()})
		}
	}

	return directories, nil
}

func (s *SCollection) retrieve(collectionName string) (Collection, error) {
	collectionPath, err := s.Storage.Main()
	if err != nil {
		return Collection{}, err
	}

	_ = fmt.Sprintf("%v", collectionPath)

	return Collection{
		Name: collectionName,
		//Paths: collectionPath,
	}, nil
}
