package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"time"
)

type DAO struct {
	log logz.Logger
	db  DBTX
	// txDepth reports how many times a transaction was started and closed
	txDepth int
	*Queries
}

func NewDAO(ctx context.Context, log logz.Logger, dsn string) (*sql.DB, *DAO, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, log.Errorw(ctx, err, "cannot open database")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, nil, log.Errorw(ctx, err, "cannot ping database")
	}

	return db, &DAO{
		log:     log.Named("dbdao"),
		db:      db,
		Queries: New(db),
	}, nil
}

func (q *DAO) Ping(ctx context.Context) error {
	return q.db.(*sql.DB).PingContext(ctx)
}

type beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func (q *DAO) BeginTx(ctx context.Context) (*DAO, func() error, error) {
	tx, err := q.db.(beginner).BeginTx(ctx, nil)
	if err != nil {
		return nil, func() error { return nil }, q.log.Errorw(ctx, err, "error while starting transaction")
	}

	txLog := q.log.With("tx", q.txDepth+1)
	return &DAO{
			log:     txLog,
			db:      q.db,
			txDepth: q.txDepth + 1,
			Queries: q.WithTx(tx),
		}, func() error {
			if errors.Is(tx.Rollback(), sql.ErrTxDone) {
				// This callback can be called as deferred, regardless if tx was committed or not, thus we should
				// silence sql.ErrTxDone error.
				q.Queries.db = q.db
				q.log.Debugw(ctx, "transaction committed")
				return nil
			}
			return q.log.Errorw(ctx, err, "error while rolling back transaction")
		}, nil
}

func (q *DAO) Commit(ctx context.Context) error {
	err := q.Queries.db.(*sql.Tx).Commit()
	q.Queries.db = q.db
	return q.log.Errorw(ctx, err, "error while committing transaction")
}

func NilStr(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
