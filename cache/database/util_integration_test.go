package database

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_integration_ValidateDbAbsPath(t *testing.T) {
	abs, err := filepath.Abs("../../test/fixture/db/VACUUM.cache")
	if err != nil {
		t.Fatal(err)
	}

	ok, err := ValidateDbAbsPath(abs)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, ok)
}

func Test_integration_Error_ValidateDbAbsPath_fake_path(t *testing.T) {
	abs, err := filepath.Abs("../../test/fixture/db/none_existing_file.cache")
	if err != nil {
		t.Fatal(err)
	}

	ok, err := ValidateDbAbsPath(abs)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to stat file")
	assert.False(t, ok)
}
