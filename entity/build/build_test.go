package build

import (
	"errors"
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
		//
		// Error
		//
		{
			name:                                   "Error with entity URI",
			entityUri:                              "github.com",
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Have:  false,
				Exist: false,
			},
			hasError: true,
		},
		{
			name:                                   "Error from entity version List",
			entityUri:                              "github.com/AmadlaOrg/Entity",
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           errors.New("error from List"),
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				RepoUrl: "https://github.com/AmadlaOrg/Entity",
				Have:    false,
				Exist:   false,
			},
			hasError: true,
		},
		{
			name:                                   "Error from entity pseudo version generator",
			entityUri:                              "github.com/AmadlaOrg/Entity",
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: errors.New("error from pseudo version generator"),
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				RepoUrl: "https://github.com/AmadlaOrg/Entity",
				Have:    false,
				Exist:   false,
			},
			hasError: true,
		},
		{
			name:                                   "Error from entity Latest",
			entityUri:                              "github.com/AmadlaOrg/Entity",
			internalEntityVersionList:              []string{""},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         errors.New("error from pseudo version generator"),
			expectEntity: entity.Entity{
				RepoUrl: "https://github.com/AmadlaOrg/Entity",
				Have:    false,
				Exist:   false,
			},
			hasError: true,
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

func TestMetaFromRemoteWithVersion(t *testing.T) {
	tests := []struct {
		name                                   string
		entityUri                              string
		internalEntityVersionExtract           string
		internalEntityVersionExtractErr        error
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
			entityUri:                              "github.com/AmadlaOrg/Entity@v1.0.0",
			internalEntityVersionExtract:           "v1.0.0",
			internalEntityVersionExtractErr:        nil,
			internalEntityVersionList:              []string{"v1.0.0", "v2.0.0", "v3.0.0"},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "v1.0.0",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Entity:          "github.com/AmadlaOrg/Entity@v1.0.0",
				Name:            "Entity",
				RepoUrl:         "https://github.com/AmadlaOrg/Entity",
				Origin:          "github.com/AmadlaOrg/",
				Version:         "v1.0.0",
				IsPseudoVersion: false,
				Have:            false,
				Exist:           false,
			},
			hasError: false,
		},
		/*{
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
		},*/
		//
		// Error
		//
		{
			name:                                   "Error with entity version extract",
			entityUri:                              "github.com/AmadlaOrg/Entity@v1.0.0",
			internalEntityVersionExtract:           "",
			internalEntityVersionExtractErr:        errors.New("error from Extract"),
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				IsPseudoVersion: false,
				Have:            false,
				Exist:           false,
			},
			hasError: true,
		},
		{
			name:                                   "Error extracting repo url",
			entityUri:                              "github.com/@v1.0.0", // TODO: Without the `/` the ExtractRepoUrl does not throw an error
			internalEntityVersionExtract:           "",
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Have:  false,
				Exist: false,
			},
			hasError: true,
		},
		{
			name:                                   "Error from entity version List",
			entityUri:                              "github.com/AmadlaOrg/Entity@latest",
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           errors.New("error from List"),
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				RepoUrl:         "https://github.com/AmadlaOrg/Entitylatest", // TODO: Correct the repo URL
				IsPseudoVersion: false,
				Have:            false,
				Exist:           false,
			},
			hasError: true,
		},
		{
			name:                                   "Error from entity pseudo version generator",
			entityUri:                              "github.com/AmadlaOrg/Entity@latest",
			internalEntityVersionExtract:           "latest",
			internalEntityVersionExtractErr:        nil,
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: errors.New("error from pseudo version generator"),
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				RepoUrl: "https://github.com/AmadlaOrg/Entity",
				Have:    false,
				Exist:   false,
			},
			hasError: true,
		},
		{
			name:                                   "Error from entity Latest",
			entityUri:                              "github.com/AmadlaOrg/Entity@latest",
			internalEntityVersionExtract:           "latest",
			internalEntityVersionExtractErr:        nil,
			internalEntityVersionList:              []string{""},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         errors.New("error from pseudo version generator"),
			expectEntity: entity.Entity{
				RepoUrl: "https://github.com/AmadlaOrg/Entity",
				Have:    false,
				Exist:   false,
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVersion := version.MockVersion{}
			mockVersion.EXPECT().Extract(mock.Anything).Return(
				tt.internalEntityVersionExtract,
				tt.internalEntityVersionExtractErr)
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

			entityMeta, err := mockBuilder.metaFromRemoteWithVersion(tt.entityUri)
			if tt.hasError {
				assert.Error(t, err)
			}

			if !reflect.DeepEqual(entityMeta, tt.expectEntity) {
				t.Errorf("got %v, want %v", entityMeta, tt.expectEntity)
			}
		})
	}
}
