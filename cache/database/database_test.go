package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitialize(t *testing.T) {
	databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	isInitialized := databaseService.IsInitialized()
	assert.True(t, isInitialized)

	err = databaseService.Close()
	assert.NoError(t, err)
}

func TestClose(t *testing.T) {}

func TestIsInitialized(t *testing.T) {}

func TestCreateTable(t *testing.T) {
	databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	err = databaseService.CreateTable(Table{Name: "test"})
	assert.NoError(t, err)

	err = databaseService.Close()
	assert.NoError(t, err)
}

func TestInsert(t *testing.T) {
	databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	err = databaseService.Insert(Table{Name: "test"}, []string{"Joe", "Jack"})
	assert.NoError(t, err)

	err = databaseService.Close()
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {}

func TestUpdate(t *testing.T) {}

func TestSelect(t *testing.T) {}

func TestDropTable(t *testing.T) {
	databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	err = databaseService.DropTable(Table{Name: "test"})
	assert.NoError(t, err)

	err = databaseService.Close()
	assert.NoError(t, err)
}
