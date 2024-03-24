package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ctxKey struct{}

type Queryer interface {
	sqlx.ExecerContext
	sqlx.QueryerContext

	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type DB struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *DB {
	return &DB{
		db: db,
	}
}

func (db *DB) WithTx(ctx context.Context, txFn func(ctx context.Context) error) (err error) {
	if tx := ctx.Value(ctxKey{}); tx != nil { // flatten transactions
		return txFn(ctx)
	}

	tx, err := db.db.BeginTxx(ctx, nil)
	if err != nil {
		return
	}
	ctx = context.WithValue(ctx, ctxKey{}, tx)

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFn(ctx)
	return
}

func (db *DB) Get(ctx context.Context) Queryer {
	if tx := ctx.Value(ctxKey{}); tx != nil {
		return tx.(*sqlx.Tx)
	}

	return db.db
}
