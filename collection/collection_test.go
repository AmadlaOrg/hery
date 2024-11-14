package collection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelect(t *testing.T) {}

func TestCreate(t *testing.T) {}

func TestRemove(t *testing.T) {}

func TestExists(t *testing.T) {}

func TestList(t *testing.T) {
	collectionService := NewCollectionService()
	list, err := collectionService.List()
	if err != nil {
		return
	}
	assert.NotNil(t, list)
}
