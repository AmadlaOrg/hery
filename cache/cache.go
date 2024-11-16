package cache

import (
	"errors"
	"github.com/AmadlaOrg/hery/cache/database"
	"github.com/AmadlaOrg/hery/cache/parser"
	"github.com/AmadlaOrg/hery/entity"
)

// ICache
type ICache interface {
	Open() error
	Close() error
	AddEntity() error
	InsertInEntity() error
	SelectEntity() (entity.Entity, error)
}

// SCache
type SCache struct {
	Database database.IDatabase
	Parser   parser.IParser
}

// Open the cache (connects to the SQLite3 file database)
func (s *SCache) Open() error {
	err := s.Database.Initialize()
	if err != nil {
		return err
	}

	return nil
}

// Close the cache (closes the SQLite3 file database)
func (s *SCache) Close() error {
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

// AddEntity
func (s *SCache) AddEntity() error {
	return nil
}

// InsertInEntity
func (s *SCache) InsertInEntity() error {
	return nil
}

// SelectEntity
func (s *SCache) SelectEntity() (entity.Entity, error) {
	return entity.Entity{}, nil
}
