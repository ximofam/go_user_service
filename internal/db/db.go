package db

import (
	"context"
	"database/sql"

	"github.com/ximofam/user-service/internal/utils"
)

type SQLScanner interface {
	Scan(dest ...any) error
}

type Querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type DBProvider struct {
	db *sql.DB
}

func NewDBProvider(db *sql.DB) *DBProvider {
	return &DBProvider{db: db}
}

func (dp *DBProvider) GetQuerier(ctx context.Context) Querier {
	tx := utils.ContextGetTx(ctx)
	if tx != nil {
		return tx
	}

	return dp.db
}
