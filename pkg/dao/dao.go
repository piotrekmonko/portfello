package dao

import (
	"database/sql"
	"fmt"
)

func NewDAO(dsn string) (*sql.DB, *Queries, error) {
	db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
	if err != nil {
		return nil, nil, fmt.Errorf("cannot reach database: %w", err)
	}

	return db, New(db), nil
}
