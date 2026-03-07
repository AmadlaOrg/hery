package parser

import (
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/AmadlaOrg/hery/entity/schema"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestParseEntity(t *testing.T) {
	mockDb := &database.MockCacheDatabase{}
	mockDb.EXPECT().IsInitialized().Return(true)
	mockDb.EXPECT().Insert(mock.Anything).Return()
	mockDb.EXPECT().Apply().Return(nil)
	parserService := NewParserService(mockDb)

	e := entity.Entity{
		Id:              uuid.MustParse("c0fdd76d-a5b5-4f35-8784-e6238d6933ab"),
		Uri:             "github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0",
		Name:            "WebServer",
		RepoUrl:         "https://github.com/AmadlaOrg/EntityApplication",
		Origin:          "github.com/AmadlaOrg/EntityApplication",
		Version:         "v1.0.0",
		IsLatestVersion: true,
		IsPseudoVersion: false,
		AbsPath:         "/home/user/.cache/hery/entity/github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0",
		Have:            true,
		Hash:            "",
		Exist:           true,
		Schema:          &schema.Schema{},
		Content: entity.Content{
			Type: "github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0",
			Self: "c0fdd76d-a5b5-4f35-8784-e6238d6933ab",
			Meta: map[string]any{
				"_type": "github.com/AmadlaOrg/Entity@latest",
				"_body": map[string]any{
					"name":        "WebServer",
					"description": "",
					"tags": []string{
						"server",
						"web",
						"service",
					},
				},
			},
			Body: map[string]any{
				"server_name": "localhost",
			},
		},
	}

	dbTable, err := parserService.Entity(&e)
	assert.NoError(t, err)
	assert.Len(t, dbTable, 1)
	assert.Equal(t, "entities", dbTable[0].Name)
	assert.Len(t, dbTable[0].Rows, 1)

	row := dbTable[0].Rows[0]
	assert.Equal(t, "github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0", row["entity_type"])
	assert.Equal(t, "c0fdd76d-a5b5-4f35-8784-e6238d6933ab", row["entity_self"])
	assert.NotEmpty(t, row["meta_json"])
	assert.NotEmpty(t, row["body_json"])
	assert.NotEmpty(t, row["merged_json"])
}

func TestEntityToTableName(t *testing.T) {
	parserService := NewParserService(nil)
	tableName := parserService.EntityToTableName("github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0")
	assert.Equal(t, "github_com_AmadlaOrg_EntityApplication_WebServer_v1_0_0", tableName)
}
