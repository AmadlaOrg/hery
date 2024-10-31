package get

import (
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name                string
		inputCollectionName string
		inputStoragePaths   *storage.AbsPaths
		inputEntitiesMeta   []entity.Entity
		expectedIds         []string
		hasError            bool
	}{
		{
			name:                "Valid",
			inputCollectionName: "amadla",
			inputStoragePaths: &storage.AbsPaths{
				Storage:    "/home/user/.hery/",
				Catalog:    "/home/user/.hery/collection/",
				Collection: "/home/user/.hery/collection/amadla/",
				Entities:   "/home/user/.hery/collection/amadla/entity/",
				Cache:      "/home/user/.hery/collection/amadla/amadla.cache",
			},
			inputEntitiesMeta: []entity.Entity{
				{
					Id:              "97d4b783-f448-483c-8111-380d6082ae1c",
					Entity:          "github.com/AmadlaOrg/Entity@v0.0.0-20240924093300-abcd1234efgh",
					Name:            "Entity",
					RepoUrl:         "https://github.com/AmadlaOrg/Entity",
					Origin:          "github.com/AmadlaOrg/Entity",
					Version:         "latest",
					LatestVersion:   true,
					IsPseudoVersion: true,
					AbsPath:         "/home/user/.hery/collection/amadla/entity/github.com/AmadlaOrg/Entity@v.0.0.0-20240924093300-abcd1234efgh",
					Have:            true,
					Hash:            "",
					Exist:           true,
					Schema:          &jsonschema.Schema{},
					Config:          nil,
				},
				{
					Id:              "12c4b793-d458-756f-8151-740d6082ae1f",
					Entity:          "github.com/AmadlaOrg/Entity@v0.0.0-20230924093300-abcd1234efgh",
					Name:            "Entity",
					RepoUrl:         "https://github.com/AmadlaOrg/Entity",
					Origin:          "github.com/AmadlaOrg/Entity",
					Version:         "latest",
					LatestVersion:   false,
					IsPseudoVersion: true,
					AbsPath:         "/home/user/.hery/collection/amadla/entity/github.com/AmadlaOrg/Entity@v.0.0.0-20230924093300-abcd1234efgh",
					Have:            true,
					Hash:            "",
					Exist:           true,
					Schema:          &jsonschema.Schema{},
					Config:          nil,
				},
				{
					Id:              "98d4b682-c758-943c-8911-560d9022ae3c",
					Entity:          "github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion@v2.1.0",
					Name:            "QAFixturesEntityMultipleTagVersion",
					RepoUrl:         "https://github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion",
					Origin:          "github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion",
					Version:         "v2.1.0",
					LatestVersion:   true,
					IsPseudoVersion: false,
					AbsPath:         "/home/user/.hery/amadla/entity/github.com/AmadlaOrg/QAFixturesEntityMultipleTagVersion@v2.1.0",
					Have:            true,
					Hash:            "",
					Exist:           true,
					Schema:          &jsonschema.Schema{},
					Config:          nil,
				},
			},
			expectedIds: []string{},
			hasError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getService := NewGetService()
			err := getService.download(tt.inputCollectionName, tt.inputStoragePaths, tt.inputEntitiesMeta)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedIds, tt.inputCollectionName)
		})
	}
}
