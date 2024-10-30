package build

import (
	"errors"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/validation"
	"github.com/AmadlaOrg/hery/entity/version"
	versionValidationPkg "github.com/AmadlaOrg/hery/entity/version/validation"
	"github.com/AmadlaOrg/hery/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestMeta(t *testing.T) {
	// Mocking UUID generation for consistent results
	originalUUIDNew := uuidNew
	defer func() { uuidNew = originalUUIDNew }() // Restore original function after test

	uuidNew = func() uuid.UUID {
		return uuid.MustParse("dca736d3-26c4-46b2-be5a-dfbdc09cff6d")
	}

	tests := []struct {
		name                 string
		inputPaths           storage.AbsPaths
		inputEntityUri       string
		internalEntityDir    string
		internalEntityDirErr error
		mockValidation       func(entityValidation *validation.MockEntityValidation)
		mockEntityVersion    func(entityVersion *version.MockEntityVersion)
		mockEntityVersionVal func(versionValidation *versionValidationPkg.MockEntityVersionValidation)
		expectEntity         entity.Entity
		hasError             bool
	}{
		{
			name: "Valid Entity URI With Version",
			inputPaths: storage.AbsPaths{
				Entities: "/home/user/.hery/amadla/entity/",
			},
			inputEntityUri:    "github.com/example/entity@v1.0.0",
			internalEntityDir: "/home/user/.hery/amadla/entity/github.com/example/entity@v1.0.0",
			mockValidation: func(mockValidation *validation.MockEntityValidation) {
				mockValidation.EXPECT().EntityUri("github.com/example/entity@v1.0.0").Return(true)
			},
			mockEntityVersion: func(mockEntityVersion *version.MockEntityVersion) {
				mockEntityVersion.EXPECT().Extract("github.com/example/entity@v1.0.0").Return("v1.0.0", nil)
			},
			mockEntityVersionVal: func(mockEntityVersionVal *versionValidationPkg.MockEntityVersionValidation) {
				// No specific mocks needed for validation in this case
			},
			expectEntity: entity.Entity{
				Id:              "dca736d3-26c4-46b2-be5a-dfbdc09cff6d",
				Entity:          "github.com/example/entity@v1.0.0",
				Name:            "entity",
				RepoUrl:         "https://github.com/example/entity",
				Origin:          "github.com/example/",
				Version:         "v1.0.0",
				LatestVersion:   false,
				IsPseudoVersion: false,
				AbsPath:         "/home/user/.hery/amadla/entity/github.com/example/entity@v1.0.0",
				Have:            true,
				Hash:            "",
				Exist:           true,
				Schema:          nil,
				Config:          nil,
			},
			hasError: false,
		},
		{
			name: "Valid Entity URI With Latest Version",
			inputPaths: storage.AbsPaths{
				Entities: "/home/user/.hery/amadla/entity/",
			},
			inputEntityUri:       "github.com/example/entity@latest",
			internalEntityDir:    "",
			internalEntityDirErr: entity.ErrorNotFound, // Not found since `latest` should never be used in directory name (only the latest static version or pseudo version)
			mockValidation: func(mockValidation *validation.MockEntityValidation) {
				mockValidation.EXPECT().EntityUri("github.com/example/entity@latest").Return(true)
			},
			mockEntityVersion: func(mockEntityVersion *version.MockEntityVersion) {
				mockEntityVersion.EXPECT().Extract("github.com/example/entity@latest").Return("latest", nil)
				mockEntityVersion.EXPECT().List("https://github.com/example/entity").Return([]string{"v1.0.0", "v1.0.1"}, nil)
				mockEntityVersion.EXPECT().Latest([]string{"v1.0.0", "v1.0.1"}).Return("v1.0.1", nil)
			},
			mockEntityVersionVal: func(mockEntityVersionVal *versionValidationPkg.MockEntityVersionValidation) {
				// No specific mocks needed for validation in this case
				//mockEntityVersionVal.EXPECT().
			},
			expectEntity: entity.Entity{
				Id:              "dca736d3-26c4-46b2-be5a-dfbdc09cff6d",
				Entity:          "github.com/example/entity@v1.0.1",
				Name:            "entity",
				RepoUrl:         "https://github.com/example/entity",
				Origin:          "github.com/example/",
				Version:         "v1.0.1",
				LatestVersion:   true,
				IsPseudoVersion: false,
				AbsPath:         "/home/user/.hery/amadla/entity/github.com/example/entity@v1.0.1",
				Have:            false,
				Hash:            "",
				Exist:           true,
				Schema:          nil,
				Config:          nil,
			},
			hasError: false,
		},
		//
		// Error
		//
		{
			name:              "Invalid Entity URI",
			inputPaths:        storage.AbsPaths{Entities: "testdata"},
			inputEntityUri:    "invalid_uri",
			internalEntityDir: "",
			mockValidation: func(mockValidation *validation.MockEntityValidation) {
				mockValidation.EXPECT().EntityUri("invalid_uri").Return(false)
			},
			mockEntityVersion:    func(mockEntityVersion *version.MockEntityVersion) {},
			mockEntityVersionVal: func(mockEntityVersionVal *versionValidationPkg.MockEntityVersionValidation) {},
			expectEntity: entity.Entity{
				Have:  false,
				Exist: false,
			},
			hasError: true,
		},
		{
			name:              "Error from metaFromRemoteWithoutVersion",
			inputPaths:        storage.AbsPaths{Entities: "testdata"},
			inputEntityUri:    "https://github.com/example/entity",
			internalEntityDir: "",
			mockValidation: func(mockValidation *validation.MockEntityValidation) {
				mockValidation.EXPECT().EntityUri("https://github.com/example/entity").Return(true)
			},
			mockEntityVersion: func(mockEntityVersion *version.MockEntityVersion) {
				// Ensure that List is mocked even though we expect an error
				mockEntityVersion.EXPECT().List("https://github.com/example/entity").Return(nil, errors.New("simulated List error"))
			},
			mockEntityVersionVal: func(mockEntityVersionVal *versionValidationPkg.MockEntityVersionValidation) {
				// No specific mocks needed for validation in this case
			},
			expectEntity: entity.Entity{
				RepoUrl:         "https://github.com/example/entity",
				IsPseudoVersion: false,
				Have:            false,
				Exist:           false,
			},
			hasError: true,
		},
		// FIXME:
		/*{
			name:                 "Error FindEntityDir with MultipleFoundError",
			inputPaths:           storage.AbsPaths{Entities: "testdata"},
			inputEntityUri:       "https://github.com/example/entity",
			internalEntityDir:    "",
			internalEntityDirErr: errtypes.MultipleFoundError,
			mockValidation: func(mockValidation *validation.MockEntityValidation) {
				mockValidation.EXPECT().EntityUri("https://github.com/example/entity").Return(true)
			},
			mockEntityVersion: func(mockEntityVersion *version.MockEntityVersion) {
				mockEntityVersion.EXPECT().List("https://github.com/example/entity").Return([]string{"v1.0.0"}, nil)
				mockEntityVersion.EXPECT().Latest([]string{"v1.0.0"}).Return("v1.0.0", nil)
			},
			mockEntityVersionVal: func(mockEntityVersionVal *versionValidationPkg.MockEntityVersionValidation) {
				// No specific mocks needed for validation in this case
			},
			expectEntity: entity.Entity{
				Id:              "dca736d3-26c4-46b2-be5a-dfbdc09cff6d",
				Entity:          "https://github.com/example/entity@v1.0.0",
				Name:            "entity",
				RepoUrl:         "https://github.com/example/entity",
				Origin:          "https://github.com/example/",
				Version:         "v1.0.0",
				IsPseudoVersion: false,
				AbsPath:         "testdata/https:/github.com/example/entity@v1.0.0", // TODO:
				Have:            false,                                              // TODO: Double check
				Exist:           true,
			},
			hasError: true,
		},
		{
			name:                 "Error FindEntityDir with any other error",
			inputPaths:           storage.AbsPaths{Entities: "testdata"},
			inputEntityUri:       "https://github.com/example/entity",
			internalEntityDir:    "",
			internalEntityDirErr: errors.New("random any other error"),
			mockValidation: func(mockValidation *validation.MockEntityValidation) {
				mockValidation.EXPECT().EntityUri("https://github.com/example/entity").Return(true)
			},
			mockEntityVersion: func(mockEntityVersion *version.MockEntityVersion) {
				mockEntityVersion.EXPECT().List("https://github.com/example/entity").Return([]string{"v1.0.0"}, nil)
				mockEntityVersion.EXPECT().Latest([]string{"v1.0.0"}).Return("v1.0.0", nil)
			},
			mockEntityVersionVal: func(mockEntityVersionVal *versionValidationPkg.MockEntityVersionValidation) {
				// No specific mocks needed for validation in this case
			},
			expectEntity: entity.Entity{
				Id:              "dca736d3-26c4-46b2-be5a-dfbdc09cff6d",
				Entity:          "https://github.com/example/entity@v1.0.0",
				Name:            "entity",
				RepoUrl:         "https://github.com/example/entity",
				Origin:          "https://github.com/example/",
				Version:         "v1.0.0",
				IsPseudoVersion: false,
				AbsPath:         "testdata/https:/github.com/example/entity@v1.0.0", // TODO:
				Have:            false,                                              // TODO: Double check
				Exist:           true,
			},
			hasError: true,
		},*/
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockEntity := entity.MockEntity{}
			mockEntity.EXPECT().FindEntityDir(mock.Anything, mock.Anything).Return(
				test.internalEntityDir, test.internalEntityDirErr)

			mockValidation := validation.NewMockEntityValidation(t)
			test.mockValidation(mockValidation)

			mockEntityVersion := version.NewMockEntityVersion(t)
			test.mockEntityVersion(mockEntityVersion)

			mockEntityVersionVal := versionValidationPkg.NewMockEntityVersionValidation(t)
			test.mockEntityVersionVal(mockEntityVersionVal)

			mockBuilder := SBuild{
				Entity:                  &mockEntity,
				EntityValidation:        mockValidation,
				EntityVersion:           mockEntityVersion,
				EntityVersionValidation: mockEntityVersionVal,
			}

			metaFromRemote, err := mockBuilder.Meta(test.inputPaths, test.inputEntityUri)
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
}

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
		inputEntityUri                         string
		inputEntityVersion                     string
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
			inputEntityUri:                         "github.com/AmadlaOrg/Entity@v1.0.0",
			inputEntityVersion:                     "v1.0.0",
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
		{
			name:                                   "Valid meta using pseudo version",
			inputEntityUri:                         "github.com/AmadlaOrg/Entity@latest",
			inputEntityVersion:                     "latest",
			internalEntityVersionList:              []string{},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "v0.0.0-20240823005443-9b4947da3948",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Id:              "",
				Entity:          "github.com/AmadlaOrg/Entity@v0.0.0-20240823005443-9b4947da3948",
				Name:            "Entity",
				RepoUrl:         "https://github.com/AmadlaOrg/Entity",
				Origin:          "github.com/AmadlaOrg/",
				Version:         "v0.0.0-20240823005443-9b4947da3948",
				LatestVersion:   true,
				IsPseudoVersion: true,
				AbsPath:         "",
				Have:            false,
				Hash:            "",
				Exist:           false,
				Schema:          nil,
				Config:          nil,
			},
			hasError: false,
		},
		{
			name:                                   "Valid meta using latest version with latest tags",
			inputEntityUri:                         "github.com/AmadlaOrg/Entity@latest",
			inputEntityVersion:                     "latest",
			internalEntityVersionList:              []string{"v1.0.0", "v2.0.0", "v3.0.0"},
			internalEntityVersionListErr:           nil,
			internalEntityVersionGeneratePseudo:    "",
			internalEntityVersionGeneratePseudoErr: nil,
			internalEntityVersionLatest:            "v3.0.0",
			internalEntityVersionLatestErr:         nil,
			expectEntity: entity.Entity{
				Entity:          "github.com/AmadlaOrg/Entity@v3.0.0",
				Name:            "Entity",
				RepoUrl:         "https://github.com/AmadlaOrg/Entity",
				Origin:          "github.com/AmadlaOrg/",
				Version:         "v3.0.0",
				LatestVersion:   true,
				IsPseudoVersion: false,
				AbsPath:         "",
				Have:            false,
				Hash:            "",
				Exist:           false,
				Schema:          nil,
				Config:          nil,
			},
			hasError: false,
		},
		// FIXME:
		/*{
			name:               "Valid meta using no set version with version list",
			inputEntityUri:     "github.com/AmadlaOrg/Entity",
			inputEntityVersion: "",
			//internalEntityVersionExtractErr:        version.ErrorExtractNoVersionFound,
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
			name:               "Valid meta using no set version with pseudo version",
			inputEntityUri:     "github.com/AmadlaOrg/Entity",
			inputEntityVersion: "",
			//internalEntityVersionExtractErr:        version.ErrorExtractNoVersionFound,
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
			name:               "Error with entity version extract",
			inputEntityUri:     "github.com/AmadlaOrg/Entity@v1.0.0",
			inputEntityVersion: "",
			//internalEntityVersionExtractErr:        errors.New("error from Extract"),
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
		},*/
		{
			name:                                   "Error extracting repo url",
			inputEntityUri:                         "github.com/@v1.0.0", // TODO: Without the `/` the ExtractRepoUrl does not throw an error
			inputEntityVersion:                     "",
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
			inputEntityUri:                         "github.com/AmadlaOrg/Entity@latest",
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
			name:               "Error from entity pseudo version generator",
			inputEntityUri:     "github.com/AmadlaOrg/Entity@latest",
			inputEntityVersion: "latest",
			//internalEntityVersionExtractErr:        nil,
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
			name:               "Error from entity Latest",
			inputEntityUri:     "github.com/AmadlaOrg/Entity@latest",
			inputEntityVersion: "latest",
			//internalEntityVersionExtractErr:        nil,
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
		{
			name:               "Invalid pass version in the entity URI",
			inputEntityUri:     "github.com/AmadlaOrg/Entity@v0.0",
			inputEntityVersion: "v0.0",
			//internalEntityVersionExtractErr:        nil,
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVersion := &version.MockEntityVersion{}
			/*mockVersion.EXPECT().Extract(mock.Anything).Return(
			tt.internalEntityVersionExtract,
			tt.internalEntityVersionExtractErr)*/
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

			mockBuilder := SBuild{
				EntityVersion:           mockVersion,
				EntityValidation:        mockEntityValidation,
				Entity:                  mockEntity,
				EntityVersionValidation: &versionValidationPkg.SValidation{},
			}

			entityMeta, err := mockBuilder.metaFromRemoteWithVersion(tt.inputEntityUri, tt.inputEntityVersion)
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

func TestMetaFromLocalWithVersion_Error_FromExtractRepoUrl(t *testing.T) {
	builder := SBuild{}
	_, err := builder.metaFromLocalWithVersion("https://github.com", "v1.0.0")

	assert.Error(t, err)
}

func TestConstructOrigin(t *testing.T) {
	mockBuilder := SBuild{}

	tests := []struct {
		name           string
		inputEntityUri string
		inputName      string
		inputVersion   string
		expected       string
	}{
		{
			name:           "Valid origin",
			inputEntityUri: "github.com/AmadlaOrg/Entity",
			inputName:      "github.com/AmadlaOrg/Entity@latest",
			inputVersion:   "latest",
			expected:       "github.com/AmadlaOrg/Entity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mockBuilder.constructOrigin(tt.inputEntityUri, tt.inputName, tt.inputVersion)
			assert.Equal(t, tt.expected, got)
		})
	}
}
