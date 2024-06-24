package collection

import "os"

func Get(storagePath string) (Collection, error) {
	if storagePath == "" {
		storagePath = "~/.hery"
	}

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		return Collection{}, err
	}

	if collectionList, err := List(storagePath); err != nil {
		return Collection{}, err
	}

}

func List(storagePath string) (Collections, error) {
	if storagePath == "" {

	}
}
