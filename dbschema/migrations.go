package dbschema

import (
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func NewMigrator(dsn string) (*migrate.Migrate, error) {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return nil, fmt.Errorf("cannot read migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot reach database, is DB DSN valid?: %w", err)
	}

	return m, nil
}
