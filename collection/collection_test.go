package collection

import (
	gitConfig "github.com/AmadlaOrg/LibraryUtils/git/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelect(t *testing.T) {
	collectionService := NewCollectionService(&gitConfig.Config{})
	selectCollection := collectionService.Select("amadla")
	assert.NotNil(t, selectCollection)
}

func TestCreate(t *testing.T) {}

func TestRemove(t *testing.T) {}

func TestExists(t *testing.T) {}

func TestList(t *testing.T) {
	collectionService := NewCollectionService(&gitConfig.Config{})
	list, err := collectionService.List()
	if err != nil {
		return
	}
	assert.NotNil(t, list)
}
