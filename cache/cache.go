package cache

import (
	"errors"
	"fmt"
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/cache/parser"
	"github.com/AmadlaOrg/hery/entity"
)

// Cache defines the cache interface for entity storage.
type Cache interface {
	Open() error
	Close() error
	AddEntity(entity *entity.Entity) error
	SelectEntity(uri string) (*entity.Entity, error)
}

// cacheImpl
type cacheImpl struct {
	Database database.Database
	Parser   parser.Parser
}

// Open the cache (connects to the SQLite3 file database)
func (s *cacheImpl) Open() error {
	err := s.Database.Initialize()
	if err != nil {
		return err
	}

	return nil
}

// Close the cache (closes the SQLite3 file database)
func (s *cacheImpl) Close() error {
	// 1. Check if the DB is initialized
	if s.Database.IsInitialized() {

		// 2. If it was initialized then it closes the connection and if there are any errors it returns them
		if err := s.Database.Close(); err != nil {
			return errors.Join(errors.New(database.ErrorClosingDatabase), err)
		}
	} else {
		return errors.New(database.ErrorDatabaseNotInitialized)
	}

	return nil
}

// AddEntity parses an entity and inserts it into the SQLite cache.
func (s *cacheImpl) AddEntity(e *entity.Entity) error {
	// Ensure tables exist
	s.Database.CreateTable()
	if err := s.Database.Apply(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// Parse and insert entity data
	_, err := s.Parser.Entity(e)
	if err != nil {
		return fmt.Errorf("failed to parse entity for cache: %w", err)
	}

	return nil
}

// SelectEntity queries an entity from the cache by URI.
func (s *cacheImpl) SelectEntity(uri string) (*entity.Entity, error) {
	rows, err := s.Database.QueryRows(
		"SELECT entity_type, entity_extends, uri, name, repo_url, origin, version, "+
			"is_latest_version, is_pseudo_version, abs_path, have, hash, exist, schema_json, "+
			"meta_json, body_json, merged_json FROM entities WHERE uri = ?", uri)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("entity not found: %s", uri)
	}

	row := rows[0]
	e := &entity.Entity{
		Uri:     stringVal(row["uri"]),
		Name:    stringVal(row["name"]),
		RepoUrl: stringVal(row["repo_url"]),
		Origin:  stringVal(row["origin"]),
		Version: stringVal(row["version"]),
		AbsPath: stringVal(row["abs_path"]),
		Hash:    stringVal(row["hash"]),
		Content: entity.Content{
			Type:    stringVal(row["entity_type"]),
			Extends: stringVal(row["entity_extends"]),
		},
	}

	return e, nil
}

func stringVal(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
