package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type DAO struct {
	log *log.Logger
	db  DBTX
	*Queries
}

func NewDAO(ctx context.Context, dsn string) (*sql.DB, *DAO, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot reach database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("cannot reach database: %w", err)
	}

	return db, &DAO{
		log:     log.New(os.Stdout, "dbdao", log.LstdFlags),
		db:      db,
		Queries: New(db),
	}, nil
}

type beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func (q *DAO) BeginTx(ctx context.Context) (*DAO, func() error, error) {
	tx, err := q.db.(beginner).BeginTx(ctx, nil)
	if err != nil {
		return nil, func() error { return nil }, fmt.Errorf("cannot start transaction: %w", err)
	}

	return &DAO{
			log:     q.log,
			db:      q.db,
			Queries: q.WithTx(tx),
		}, func() error {
			if errors.Is(tx.Rollback(), sql.ErrTxDone) {
				// This callback can be called as deferred, regardless if tx was commited or not, thus we should
				// silence sql.ErrTxDone error.
				q.Queries.db = q.db
				return nil
			}
			q.log.Printf("error while committing transaction: %s", err)
			return err
		}, nil
}

func (q *DAO) Commit() error {
	err := q.Queries.db.(*sql.Tx).Commit()
	q.Queries.db = q.db
	return err
}

func NilStr(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}
