package sql

import (
	"context"
	"database/sql"
)

type contextKey string

const txKey = contextKey("tx")

func getTx(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	if ok {
		return tx
	}
	return nil
}

func withTx(ctx context.Context, pipe *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, pipe)
}

func Transitionally(
	ctx context.Context,
	db *sql.DB,
	fn func(context.Context) error,
) error {
	if getTx(ctx) != nil {
		return fn(ctx)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctx = withTx(ctx, tx)

	err = fn(ctx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func GetTx(ctx context.Context) *sql.Tx {
	tx := getTx(ctx)
	if tx != nil {
		return tx
	}

	panic("sql: no tx")
}
