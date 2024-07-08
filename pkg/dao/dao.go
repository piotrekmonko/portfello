package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/piotrekmonko/portfello/pkg/conf"
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

type DBInterface interface {
	Querier
	DB() *sql.DB
	Ping(ctx context.Context) error
	BeginTx(ctx context.Context) (DBInterface, func(), error)
	Commit(ctx context.Context) error
}

func NewDAO(ctx context.Context, log *logz.Log, c *conf.Config) (*DAO, func(), error) {
	db, err := sql.Open("postgres", c.DatabaseDSN)
	if err != nil {
		return nil, nil, log.Errorw(ctx, err, "cannot open database")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, nil, log.Errorw(ctx, err, "cannot ping database")
	}

	return &DAO{
			log:     log.Named("dbdao"),
			db:      db,
			Queries: New(db),
		}, func() {
			_ = db.Close()
		}, nil
}

func (q *DAO) DB() *sql.DB {
	return q.db.(*sql.DB)
}

func (q *DAO) Ping(ctx context.Context) error {
	return q.db.(*sql.DB).PingContext(ctx)
}

type beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func (q *DAO) BeginTx(ctx context.Context) (DBInterface, func(), error) {
	tx, err := q.db.(beginner).BeginTx(ctx, nil)
	if err != nil {
		return nil, func() {}, q.log.Errorw(ctx, err, "error while starting transaction")
	}

	txLog := q.log.With("tx", q.txDepth+1)
	return &DAO{
			log:     txLog,
			db:      q.db,
			txDepth: q.txDepth + 1,
			Queries: q.WithTx(tx),
		}, func() {
			rollbackErr := tx.Rollback()
			if errors.Is(rollbackErr, sql.ErrTxDone) {
				// This callback can be called as deferred, regardless if tx was committed or not, thus we should
				// silence sql.ErrTxDone error.
				q.Queries.db = q.db
				q.log.Debugw(ctx, "transaction committed")
				return
			}
			_ = q.log.Errorw(ctx, rollbackErr, "error while rolling back transaction")
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
