package cache

import (
	_ "github.com/mattn/go-sqlite3"
)

// ICache
type ICache interface {
	Create()
	Insert()
	Select()
}

// SCache
type SCache struct{}

// TODO: journal mode: wal
// This will help with performance
// https://stackoverflow.com/questions/57118674/go-sqlite3-with-journal-mode-wal-gives-database-is-locked-error

// Create
func (s *SCache) Create() {

}

// Insert
func (s *SCache) Insert() {

}

// Select
func (s *SCache) Select() {

}
