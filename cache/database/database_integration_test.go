package database

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Integration_Initialize(t *testing.T) {
	/*databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	isInitialized := databaseService.IsInitialized()
	assert.True(t, isInitialized)

	err = databaseService.Close()
	assert.NoError(t, err)*/
}

func Test_Integration_IsInitialized(t *testing.T) {}

func Test_Integration_CreateTable(t *testing.T) {
	uuidString := uuid.New()
	path := fmt.Sprintf("/tmp/test-%s.db", uuidString)
	databaseService := NewDatabaseService(path)
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
	// TODO: Delete DB
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

func Test_Integration_DropTable(t *testing.T) {
	/*databaseService := NewDatabaseService()
	err := databaseService.Initialize()
	assert.NoError(t, err)

	err = databaseService.DropTable(Table{Name: "test"})
	assert.NoError(t, err)

	err = databaseService.Close()
	assert.NoError(t, err)*/
}
