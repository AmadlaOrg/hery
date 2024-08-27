package build

import (
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestMetaFromRemoteWithoutVersion(t *testing.T) {
	tests := []struct {
		name                                   string
		entityUri                              string
		internalEntityVersionList              []string
		internalEntityVersionListErr           error
		internalEntityVersionGeneratePseudo    string
		internalEntityVersionGeneratePseudoErr error
		internalEntityVersionLatest            string
		internalEntityVersionLatestErr         error
		expectEntity                           entity.Entity
		hasError                               bool
	}{
		{
			name:                                   "Valid meta using multiple entity versions",
			entityUri:                              "github.com/AmadlaOrg/Entity",
			internalEntityVersionList:              []string{"v1.0.0", "v2.0.0", "v3.0.0"},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "", //"v0.0.0-20240823005443-9b4947da3948",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "v3.0.0",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Entity:          "github.com/AmadlaOrg/Entity@v3.0.0",
				Name:            "Entity",
				RepoUrl:         "https://github.com/AmadlaOrg/Entity",
				Origin:          "github.com/AmadlaOrg/",
				Version:         "v3.0.0",
				IsPseudoVersion: false,
				Have:            false,
				Exist:           false,
			},
			hasError: false,
		},
		{
			name:                                   "Valid meta using pseudo version",
			entityUri:                              "github.com/AmadlaOrg/Entity",
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "v0.0.0-20240823005443-9b4947da3948",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Entity:          "github.com/AmadlaOrg/Entity@v0.0.0-20240823005443-9b4947da3948",
				Name:            "Entity",
				RepoUrl:         "https://github.com/AmadlaOrg/Entity",
				Origin:          "github.com/AmadlaOrg/",
				Version:         "v0.0.0-20240823005443-9b4947da3948",
				IsPseudoVersion: true,
				Have:            false,
				Exist:           false,
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVersion := version.MockVersion{}
			mockVersion.EXPECT().List(mock.Anything).Return(
				tt.internalEntityVersionList,
				tt.internalEntityVersionListErr)
			mockVersion.EXPECT().GeneratePseudo(mock.Anything).Return(
				tt.internalEntityVersionGeneratePseudo,
				tt.internalEntityVersionGeneratePseudoErr)
			mockVersion.EXPECT().Latest(mock.Anything).Return(
				tt.internalEntityVersionLatest,
				tt.internalEntityVersionLatestErr)

			mockBuilder := Builder{
				EntityVersion: &mockVersion,
			}

			entityMeta, err := mockBuilder.metaFromRemoteWithoutVersion(tt.entityUri)
			if tt.hasError {
				assert.Error(t, err)
			}

			if !reflect.DeepEqual(entityMeta, tt.expectEntity) {
				t.Errorf("got %v, want %v", entityMeta, tt.expectEntity)
			}
		})
	}
}
