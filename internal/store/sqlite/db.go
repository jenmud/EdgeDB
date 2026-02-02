package sqlite

import (
	"context"
	"embed"
	"errors"

	"github.com/jenmud/edgedb/internal/store"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/golang-migrate/migrate/v4"
	migrateSQLite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed "migrations/*.sql"
var migrations embed.FS

// New creates a new Query instance with the provided database connection.
func New(dns string) *store.DB {
	q := &store.DB{
		DB: sqlx.MustConnect("sqlite", dns),
	}

	// call SetMaxOpenConns to 1 for SQLite to avoid "database is locked" errors on the original underlying DB
	q.DB.SetMaxOpenConns(1)
	return q
}

// ApplyMigrations applies database migrations from the embedded filesystem.
func ApplyMigrations(ctx context.Context, db *store.DB) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	// db.DB.DB is a bit of inheritance mess
	driver, err := migrateSQLite.WithInstance(db.DB.DB, &migrateSQLite.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
