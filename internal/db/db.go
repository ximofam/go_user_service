package db

import (
	"context"
	"database/sql"

	"github.com/ximofam/user-service/internal/utils"
)

type SQLScanner interface {
	Scan(dest ...any) error
}

type querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db}
}

func (dp *Database) Querier(ctx context.Context) querier {
	tx, ok := ctx.Value(utils.TxKey).(*sql.Tx)
	if ok {
		return tx
	}

	return dp.db
}

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type sqlTransactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *sqlTransactor {
	return &sqlTransactor{db: db}
}

func (t *sqlTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, utils.TxKey, tx)
	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
