package database

import (
	"database/sql"
	"time"
)

// ISqlDb
type ISqlDb interface {
	Begin() (*sql.Tx, error)
	Exec(query string, args ...any) (sql.Result, error)
	Close() error
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
}

// ISqlTx
type ISqlTx interface {
	Exec(query string, args ...any) (sql.Result, error)
	Rollback() error
	Commit() error
}
