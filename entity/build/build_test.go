package build

import (
	"errors"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

/*func TestMetaFromRemote(t *testing.T) {
	tests := []struct {
		name                 string
		inputPaths           storage.AbsPaths
		inputEntityUri       string
		internalEntityDir    string
		internalEntityDirErr error
		mockValidation       func(*validation.MockInterface)
		mockEntityVersion    func(*version.MockIVersion)
		mockEntityVersionVal func(*versionValidationPkg.MockVersionValidation)
		expectEntity         entity.Entity
		hasError             bool
	}{
		{
			name:              "Valid Entity URI Without Version",
			inputPaths:        storage.AbsPaths{Entities: "testdata"},
			inputEntityUri:    "https://github.com/example/entity",
			internalEntityDir: "testdata/entity_remote.txt",
			mockValidation: func(mockValidation *validation.MockInterface) {
				mockValidation.EXPECT().EntityUrl("https://github.com/example/entity").Return(true)
			},
			mockEntityVersion: func(mockEntityVersion *version.MockIVersion) {
				mockEntityVersion.EXPECT().List("https://github.com/example/entity").Return([]string{"v1.0.0"}, nil)
				mockEntityVersion.EXPECT().Latest([]string{"v1.0.0"}).Return("v1.0.0", nil)
			},
			mockEntityVersionVal: func(mockEntityVersionVal *versionValidationPkg.MockVersionValidation) {
				// No specific mocks needed for validation in this case
			},
			expectEntity: entity.Entity{
				RepoUrl:         "https://github.com/example/entity",
				Name:            "entity",
				Version:         "v1.0.0",
				Entity:          "https://github.com/example/entity@v1.0.0",
				Origin:          "",
				AbsPath:         "testdata/entity_remote.txt",
				Have:            true,
				Exist:           true,
				IsPseudoVersion: false,
			},
			hasError: false,
		},
		{
			name:              "Invalid Entity URI",
			inputPaths:        storage.AbsPaths{Entities: "testdata"},
			inputEntityUri:    "invalid_uri",
			internalEntityDir: "",
			mockValidation: func(mockValidation *validation.MockInterface) {
				mockValidation.EXPECT().EntityUrl("invalid_uri").Return(false)
			},
			mockEntityVersion:    func(mockEntityVersion *version.MockIVersion) {},
			mockEntityVersionVal: func(mockEntityVersionVal *versionValidationPkg.MockVersionValidation) {},
			expectEntity: entity.Entity{
				Have:  false,
				Exist: false,
			},
			hasError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockEntity := entity.MockEntity{}
			mockEntity.EXPECT().FindEntityDir(mock.Anything, mock.Anything).Return(
				test.internalEntityDir, test.internalEntityDirErr)

			mockValidation := validation.NewMockInterface(t)
			test.mockValidation(mockValidation)

			mockEntityVersion := version.NewMockIVersion(t)
			test.mockEntityVersion(mockEntityVersion)

			mockEntityVersionVal := versionValidationPkg.NewMockVersionValidation(t)
			test.mockEntityVersionVal(mockEntityVersionVal)

			mockBuilder := SBuild{
				Entity:                  &mockEntity,
				EntityValidation:        mockValidation,
				EntityVersion:           mockEntityVersion,
				EntityVersionValidation: mockEntityVersionVal,
			}

			metaFromRemote, err := mockBuilder.MetaFromRemote(test.inputPaths, test.inputEntityUri)
			if test.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !reflect.DeepEqual(metaFromRemote, test.expectEntity) {
				t.Errorf("expected: %v, got: %v", test.expectEntity, metaFromRemote)
			}
		})
	}
}*/

/*func TestMetaFromRemote(t *testing.T) {
	tests := []struct {
		name                 string
		inputPaths           storage.AbsPaths
		inputEntityUri       string
		internalEntityDir    string
		internalEntityDirErr error
		expectEntity         entity.Entity
		hasError             bool
	}{
		{
			name:              "",
			inputPaths:        storage.AbsPaths{},
			inputEntityUri:    "testdata/entity_remote.txt",
			internalEntityDir: "testdata/internal_remote.txt",
			expectEntity:      entity.Entity{},
			hasError:          false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockEntity := entity.MockEntity{}
			mockEntity.EXPECT().FindEntityDir(mock.Anything, mock.Anything).Return(
				test.internalEntityDir, test.internalEntityDirErr)

			mockBuilder := SBuild{
				Entity: &mockEntity,
			}
			metaFromRemote, err := mockBuilder.MetaFromRemote(test.inputPaths, test.inputEntityUri)
			if test.hasError {
				assert.Error(t, err)
			}

			if !reflect.DeepEqual(metaFromRemote, test.expectEntity) {
				t.Errorf("expected: %v, got: %v", test.expectEntity, metaFromRemote)
			}
		})
	}
}*/

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
			mockVersion := version.MockEntityVersion{}
			mockVersion.EXPECT().List(mock.Anything).Return(
				tt.internalEntityVersionList,
				tt.internalEntityVersionListErr)
			mockVersion.EXPECT().GeneratePseudo(mock.Anything).Return(
				tt.internalEntityVersionGeneratePseudo,
				tt.internalEntityVersionGeneratePseudoErr)
			mockVersion.EXPECT().Latest(mock.Anything).Return(
				tt.internalEntityVersionLatest,
				tt.internalEntityVersionLatestErr)

			mockBuilder := SBuild{
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
		// FIXME:
		/*{
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
		},*/
		{
			name:                                   "Valid meta using pseudo version",
			entityUri:                              "github.com/AmadlaOrg/Entity@latest",
			internalEntityVersionExtract:           "latest",
			internalEntityVersionExtractErr:        nil,
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "v0.0.0-20240823005443-9b4947da3948",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Entity:          "github.com/AmadlaOrg/Entity@latest",
				Name:            "Entity",
				RepoUrl:         "https://github.com/AmadlaOrg/Entity",
				Origin:          "github.com/AmadlaOrg/Entity@latest",
				Version:         "v0.0.0-20240823005443-9b4947da3948",
				IsPseudoVersion: true,
				Have:            false,
				Exist:           false,
			},
			hasError: false,
		},
		{
			name:                                   "Valid meta using latest version with latest tags",
			entityUri:                              "github.com/AmadlaOrg/Entity@latest",
			internalEntityVersionExtract:           "latest",
			internalEntityVersionExtractErr:        nil,
			internalEntityVersionList:              []string{"v1.0.0", "v2.0.0", "v3.0.0"},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "v3.0.0",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Entity:          "github.com/AmadlaOrg/Entity@latest",
				Name:            "Entity",
				RepoUrl:         "https://github.com/AmadlaOrg/Entity",
				Origin:          "github.com/AmadlaOrg/Entity@latest",
				Version:         "v3.0.0",
				IsPseudoVersion: false,
				Have:            false,
				Exist:           false,
			},
			hasError: false,
		},
		{
			name:                                   "Valid meta using no set version with version list",
			entityUri:                              "github.com/AmadlaOrg/Entity",
			internalEntityVersionExtract:           "",
			internalEntityVersionExtractErr:        version.ErrorExtractNoVersionFound,
			internalEntityVersionList:              []string{"v1.0.0", "v2.0.0", "v3.0.0"},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "v3.0.0",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Entity:          "github.com/AmadlaOrg/Entity",
				Name:            "Entity",
				RepoUrl:         "https://github.com/AmadlaOrg/Entity",
				Origin:          "github.com/AmadlaOrg/Entity",
				Version:         "v3.0.0",
				IsPseudoVersion: false,
				Have:            false,
				Exist:           false,
			},
			hasError: false,
		},
		{
			name:                                   "Valid meta using no set version with pseudo version",
			entityUri:                              "github.com/AmadlaOrg/Entity",
			internalEntityVersionExtract:           "",
			internalEntityVersionExtractErr:        version.ErrorExtractNoVersionFound,
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "v0.0.0-20240823005443-9b4947da3948",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Entity:          "github.com/AmadlaOrg/Entity",
				Name:            "Entity",
				RepoUrl:         "https://github.com/AmadlaOrg/Entity",
				Origin:          "github.com/AmadlaOrg/Entity",
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
		// FIXME:
		/*{
			name:                                   "Invalid pass version in the entity URI",
			entityUri:                              "github.com/AmadlaOrg/Entity@v0.0",
			internalEntityVersionExtract:           "v0.0",
			internalEntityVersionExtractErr:        nil,
			internalEntityVersionList:              []string{"v1.0.0", "v2.0.0", "v3.0.0"},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Entity:          "",
				Name:            "",
				RepoUrl:         "https://github.com/AmadlaOrg/Entity",
				Origin:          "",
				Version:         "",
				IsPseudoVersion: false,
				Have:            false,
				Exist:           false,
			},
			hasError: true,
		},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVersion := &version.MockEntityVersion{}
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

			mockEntityValidation := &validation.MockEntityValidation{}
			mockEntity := &entity.MockEntity{}
			mockEntityVersionValidation := &versionValidationPkg.MockEntityVersionValidation{}

			mockBuilder := SBuild{
				EntityVersion:           mockVersion,
				EntityValidation:        mockEntityValidation,
				Entity:                  mockEntity,
				EntityVersionValidation: mockEntityVersionValidation,
			}

			entityMeta, err := mockBuilder.metaFromRemoteWithVersion(tt.entityUri)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !reflect.DeepEqual(entityMeta, tt.expectEntity) {
				t.Errorf("got %v, want %v", entityMeta, tt.expectEntity)
			}
		})
	}
}
