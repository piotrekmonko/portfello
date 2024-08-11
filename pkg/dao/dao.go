package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/piotrekmonko/portfello/dbschema"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"strings"
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

func NewDAO(ctx context.Context, logz *logz.Log, c *conf.Config) (*DAO, func(), error) {
	log := logz.With("dsn", c.DatabaseDSN)
	driver, dsn, err := driverFromDSN(c.DatabaseDSN)
	if err != nil {
		return nil, nil, log.Errorw(ctx, err, "invalid database_dsn")
	}

	log.Infow(ctx, "connecting to db", "driver", driver, "addr", dsn)
	if driver == "sqlite" {
		log.Infof(ctx, "using sqlite database, applying migrations")
		migrator, err := dbschema.NewMigrator(c.DatabaseDSN)
		if err != nil {
			return nil, nil, log.Errorw(ctx, err, "cannot init migrator for in-memory db")
		}

		err = migrator.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return nil, nil, log.Errorw(ctx, err, "cannot auto-migrate in-memory db")
		}
	}

	db, err := sql.Open(driver, dsn)
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

func driverFromDSN(dsn string) (string, string, error) {
	if dsn == "" {
		return "", "", fmt.Errorf("empty")
	}

	i := strings.Index(dsn, "://")
	if i < 1 {
		return "", "", fmt.Errorf("no echema")
	}

	if dsn[0:i] == "sqlite" {
		return dsn[0:i], dsn[i+3:], nil
	}

	return dsn[0:i], dsn, nil
}
