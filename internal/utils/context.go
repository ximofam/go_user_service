package utils

import (
	"context"
	"database/sql"
)

type txKey struct{}

var activeTxKey = txKey{}

func ContextSetTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, activeTxKey, tx)
}

func ContextGetTx(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(activeTxKey).(*sql.Tx)
	if !ok {
		return nil
	}
	return tx
}

type userIDKey struct{}

var UserIDKey = userIDKey{}
