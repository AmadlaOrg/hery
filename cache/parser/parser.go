package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/entity"
)

type IParser interface {
	Entity(entity *entity.Entity) ([]database.Table, error)
	EntityToTableName(entity string) string
	DatabaseTable(data []byte) (entity.Entity, error)
	DatabaseRow(data []byte) (entity.Entity, error)
}

type SParser struct {
	database database.IDatabase
}

var (
	jsonMarshal = json.Marshal
)

// Entity parses the entity and inserts it into the entities table.
func (s *SParser) Entity(e *entity.Entity) ([]database.Table, error) {
	if !s.database.IsInitialized() {
		return nil, errors.New("database is not initialized")
	}

	entitiesTable := s.entitiesEntityToTable(e)
	s.database.Insert(entitiesTable)
	err := s.database.Apply()
	if err != nil {
		return nil, err
	}

	return []database.Table{
		entitiesTable,
	}, nil
}

func (s *SParser) entitiesEntityToTable(e *entity.Entity) database.Table {
	metaJson, _ := jsonMarshal(e.Content.Meta)
	bodyJson, _ := jsonMarshal(e.Content.Body)

	// Build the full merged document for merged_json
	merged := map[string]any{
		"_type": e.Content.Type,
	}
	if e.Content.Self != "" {
		merged["_self"] = e.Content.Self
	}
	if e.Content.Parent != "" {
		merged["_parent"] = e.Content.Parent
	}
	if e.Content.Meta != nil {
		merged["_meta"] = e.Content.Meta
	}
	if e.Content.Body != nil {
		merged["_body"] = e.Content.Body
	}
	mergedJson, _ := jsonMarshal(merged)

	var entitiesTable database.Table
	entitiesTable.Name = "entities"
	entitiesTable.Rows = []database.Row{
		{
			"entity_type":       e.Content.Type,
			"entity_self":       e.Content.Self,
			"entity_parent":     e.Content.Parent,
			"uri":               e.Uri,
			"name":              e.Name,
			"repo_url":          e.RepoUrl,
			"origin":            e.Origin,
			"version":           e.Version,
			"is_latest_version": e.IsLatestVersion,
			"is_pseudo_version": e.IsPseudoVersion,
			"abs_path":          e.AbsPath,
			"have":              e.Have,
			"hash":              e.Hash,
			"exist":             e.Exist,
			"schema_json":       e.SchemaJson,
			"meta_json":         string(metaJson),
			"body_json":         string(bodyJson),
			"merged_json":       string(mergedJson),
		},
	}

	return entitiesTable
}

// EntityToTableName converts an entity URI to a valid SQL table name.
func (s *SParser) EntityToTableName(entity string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	tableName := re.ReplaceAllString(entity, "_")
	return strings.Trim(tableName, "_")
}

// DatabaseTable parses a JSON-encoded table row into an Entity.
func (s *SParser) DatabaseTable(data []byte) (entity.Entity, error) {
	var row map[string]any
	if err := json.Unmarshal(data, &row); err != nil {
		return entity.Entity{}, fmt.Errorf("failed to unmarshal table data: %w", err)
	}
	return s.rowToEntity(row), nil
}

// DatabaseRow parses a JSON-encoded database row into an Entity.
func (s *SParser) DatabaseRow(data []byte) (entity.Entity, error) {
	var row map[string]any
	if err := json.Unmarshal(data, &row); err != nil {
		return entity.Entity{}, fmt.Errorf("failed to unmarshal row data: %w", err)
	}
	return s.rowToEntity(row), nil
}

func (s *SParser) rowToEntity(row map[string]any) entity.Entity {
	str := func(key string) string {
		if v, ok := row[key].(string); ok {
			return v
		}
		return ""
	}
	boolVal := func(key string) bool {
		switch v := row[key].(type) {
		case bool:
			return v
		case float64:
			return v != 0
		}
		return false
	}

	e := entity.Entity{
		Uri:             str("uri"),
		Name:            str("name"),
		RepoUrl:         str("repo_url"),
		Origin:          str("origin"),
		Version:         str("version"),
		IsLatestVersion: boolVal("is_latest_version"),
		IsPseudoVersion: boolVal("is_pseudo_version"),
		AbsPath:         str("abs_path"),
		Have:            boolVal("have"),
		Hash:            str("hash"),
		Exist:           boolVal("exist"),
		SchemaJson:      str("schema_json"),
		Content: entity.Content{
			Type:   str("entity_type"),
			Self:   str("entity_self"),
			Parent: str("entity_parent"),
		},
	}

	if metaStr := str("meta_json"); metaStr != "" {
		var meta map[string]any
		if err := json.Unmarshal([]byte(metaStr), &meta); err == nil {
			e.Content.Meta = meta
		}
	}
	if bodyStr := str("body_json"); bodyStr != "" {
		var body map[string]any
		if err := json.Unmarshal([]byte(bodyStr), &body); err == nil {
			e.Content.Body = body
		}
	}

	return e
}
