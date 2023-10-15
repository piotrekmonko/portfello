package dao

import (
	"context"
	"database/sql"
	"fmt"
)

func NewDAO(dsn string) (*sql.DB, *Queries, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot reach database: %w", err)
	}

	return db, New(db), nil
}

func (q *Queries) BeginTx(ctx context.Context) (Querier, error) {
	tx, err := q.db.(*sql.DB).BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot start transaction: %w", err)
	}

	return q.WithTx(tx), nil
}
