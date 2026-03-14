package database

import (
	"context"
	"database/sql"
	"time"
)

// SqlDB
type SqlDB interface {
	Begin() (*sql.Tx, error)
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	Close() error
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
	PingContext(ctx context.Context) error
}

// SqlTx
type SqlTx interface {
	Exec(query string, args ...any) (sql.Result, error)
	Rollback() error
	Commit() error
}
