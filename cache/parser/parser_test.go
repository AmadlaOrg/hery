package parser

import (
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/entity"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseEntity(t *testing.T) {
	parserService := NewParserService()

	expected := []database.Table{
		{
			Name: "Entities",
			Columns: []database.Column{
				{
					ColumnName: "Id",
					DataType:   "string",
				},
				{
					ColumnName: "Entity",
					DataType:   "string",
				},
				{
					ColumnName: "Name",
					DataType:   "string",
				},
				{
					ColumnName: "RepoUrl",
					DataType:   "string",
				},
				{
					ColumnName: "Origin",
					DataType:   "string",
				},
				{
					ColumnName: "Version",
					DataType:   "string",
				},
				{
					ColumnName: "LatestVersion",
					DataType:   "bool",
				},
				{
					ColumnName: "IsPseudoVersion",
					DataType:   "bool",
				},
				{
					ColumnName: "AbsPath",
					DataType:   "string",
				},
				{
					ColumnName: "Have",
					DataType:   "bool",
				},
				{
					ColumnName: "Hash",
					DataType:   "string",
				},
				{
					ColumnName: "Exist",
					DataType:   "bool",
				},
				{
					ColumnName: "Schema",
					DataType:   "string",
				},
				{
					ColumnName: "Content",
					DataType:   "string",
				},
			},
			Rows: []map[string]any{
				{
					"Id":              "c0fdd76d-a5b5-4f35-8784-e6238d6933ab",
					"Entity":          "github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0",
					"Name":            "WebServer",
					"RepoUrl":         "https://github.com/AmadlaOrg/EntityApplication",
					"Origin":          "github.com/AmadlaOrg/EntityApplication",
					"Version":         "v1.0.0",
					"LatestVersion":   true,
					"IsPseudoVersion": false,
					"AbsPath":         "/home/user/.hery/amadla/entity/github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0",
					"Have":            true,
					"Hash":            "",
					"Exist":           true,
					"Schema":          ``,
					"Content":         ``,
				},
			},
		},
		{
			Name: "WebServer",
			Columns: []database.Column{
				{
					ColumnName: "Id",
					DataType:   "string",
				},
				{
					ColumnName: "server_name",
					DataType:   "string",
				},
				{
					ColumnName: "listen",
					DataType:   "string",
				},
			},
			Rows: []map[string]any{},
		},
		{
			Name: "Net",
			Columns: []database.Column{
				{
					ColumnName: "Id",
					DataType:   "string",
				},
				{
					ColumnName: "server_name",
					DataType:   "string",
				},
				{
					ColumnName: "listen",
					DataType:   "string",
				},
			},
			Rows: []map[string]any{},
		},
	}
	e := entity.Entity{
		Id:              "c0fdd76d-a5b5-4f35-8784-e6238d6933ab",
		Entity:          "github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0",
		Name:            "WebServer",
		RepoUrl:         "https://github.com/AmadlaOrg/EntityApplication",
		Origin:          "github.com/AmadlaOrg/EntityApplication",
		Version:         "v1.0.0",
		LatestVersion:   true,
		IsPseudoVersion: false,
		AbsPath:         "/home/user/.hery/amadla/entity/github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0",
		Have:            true,
		Hash:            "",
		Exist:           true,
		Schema:          &jsonschema.Schema{},
		Content: entity.Content{
			Entity: "github.com/AmadlaOrg/EntityApplication/WebServer@v1.0.0",
			Id:     "c0fdd76d-a5b5-4f35-8784-e6238d6933ab",
			Meta: map[string]any{
				"_entity": "github.com/AmadlaOrg/Entity@latest",
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
				"listen": []map[string]any{
					{
						"_entity": "github.com/AmadlaOrg/EntitySystem/Net@v1.0.0",
						"_body": map[string]any{
							"ports": []string{
								"80",
								"443",
							},
						},
					},
				},
			},
		},
	}

	dbTable, err := parserService.ParseEntity(e)
	assert.NoError(t, err)
	assert.Equal(t, expected, dbTable)
}
