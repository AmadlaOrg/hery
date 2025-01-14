package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Integration_IsInitialized_is_true(t *testing.T) {
	databaseService := NewDatabaseService("/tmp/test-integration-IsInitialized-is-true.db")
	err := databaseService.Initialize()
	assert.NoError(t, err)

	ok := databaseService.IsInitialized()
	assert.True(t, ok)

	err = databaseService.Close()
	assert.NoError(t, err)

	err = databaseService.DeleteDb()
	assert.NoError(t, err)
}

func Test_Integration_CreateTable(t *testing.T) {
	databaseService := NewDatabaseService("/tmp/test-integration-create-table.db")
	err := databaseService.Initialize()
	assert.NoError(t, err)

	databaseService.CreateTable()
	err = databaseService.Apply()
	if err != nil {
		t.Errorf("Failed to apply database %v", err)
	}

	err = databaseService.Close()
	assert.NoError(t, err)

	// TODO: Validate the content

	err = databaseService.DeleteDb()
	assert.NoError(t, err)
}

func Test_Integration_Insert(t *testing.T) {
	/*databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	err = databaseService.Insert(Table{Name: "test"}, []string{"Joe", "Jack"})
	assert.NoError(t, err)

	err = databaseService.Close()
	assert.NoError(t, err)*/
}

func Test_Integration_Select(t *testing.T) {}

func Test_Integration_Delete(t *testing.T) {}
