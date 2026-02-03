package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"github.com/golang-migrate/migrate/v4"
	migrateSQLite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed "migrations/*.sql"
var migrations embed.FS

// New creates a new Query instance with the provided database connection.
func New(dns string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", dns)
	if err != nil {
		return nil, err
	}

	// call SetMaxOpenConns to 1 for SQLite to avoid "database is locked" errors on the original underlying DB
	db.SetMaxOpenConns(1)
	return db, nil
}

// ApplyMigrations applies database migrations from the embedded filesystem.
func ApplyMigrations(ctx context.Context, db *sql.DB) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}

	// db.DB.DB is a bit of inheritance mess
	driver, err := migrateSQLite.WithInstance(db, &migrateSQLite.Config{})
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
