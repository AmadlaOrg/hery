package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Integration_Initialize(t *testing.T) {
	databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	isInitialized := databaseService.IsInitialized()
	assert.True(t, isInitialized)

	err = databaseService.Close()
	assert.NoError(t, err)
}

func Test_Integration_IsInitialized(t *testing.T) {}

func Test_Integration_CreateTable(t *testing.T) {
	databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	err = databaseService.CreateTable(Table{Name: "test"})
	assert.NoError(t, err)

	err = databaseService.Close()
	assert.NoError(t, err)
}

func Test_Integration_Insert(t *testing.T) {
	databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	err = databaseService.Insert(Table{Name: "test"}, []string{"Joe", "Jack"})
	assert.NoError(t, err)

	err = databaseService.Close()
	assert.NoError(t, err)
}

func Test_Integration_Select(t *testing.T) {}

func Test_Integration_Delete(t *testing.T) {}

func Test_Integration_DropTable(t *testing.T) {
	databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	err = databaseService.DropTable(Table{Name: "test"})
	assert.NoError(t, err)

	err = databaseService.Close()
	assert.NoError(t, err)
}
