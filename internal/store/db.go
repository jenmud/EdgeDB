package store

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/jenmud/edgedb/internal/store/sqlite"
	"github.com/jmoiron/sqlx"
)

// DB extends sqlx.DB implementing additional methods used for the store..
type DB struct {
	*sqlx.DB
}

// New creates and returns a store.
// If driver is not provided, it will look for the store driver environment variable, else will default to `sqlite`
// If dsn is not provided, it will look for the store DSN environment variable, else will default to `:memory:`
func New(ctx context.Context, driver, dsn string) (*DB, error) {
	if dsn == "" {
		dsn = ":memory:"

		envDSN := os.Getenv("EDGEDB_STORE_DSN")
		if envDSN != "" {
			dsn = envDSN
		}
	}

	if driver == "" {
		driver = "sqlite"

		envDriver := strings.ToLower(os.Getenv("EDGEDB_STORE_DRIVER"))
		if envDriver != "" {
			driver = envDriver
		}
	}

	slog.SetDefault(
		slog.With(
			slog.Group(
				"store",
				slog.String("driver", driver),
				slog.String("dsn", dsn),
			),
		),
	)

	switch strings.ToLower(driver) {
	case "duckdb":
		return nil, errors.New("duckdb not store implemented")

	case "sqlite":
		db, err := sqlite.New(dsn)
		if err != nil {
			return nil, err
		}

		slog.Info("applying db migrations")
		return &DB{DB: db}, sqlite.ApplyMigrations(ctx, db.DB)
	}

	return nil, errors.New("unsupported store")
}

// Tx starts a new transaction.
func (b *DB) Tx(ctx context.Context) (*sql.Tx, error) {
	return b.DB.BeginTx(ctx, nil)
}

// Close closed the database.
func (b *DB) Close() error {
	return b.DB.Close()
}

// Nodes returns all the nodes in the store.
func (b *DB) Nodes(ctx context.Context) ([]Node, error) {
	var nodes []Node
	return nodes, errors.New("not implemented")
}

// Edites returns all the edges in the store.
func (b *DB) Edges(ctx context.Context) ([]Edge, error) {
	var edges []Edge
	return edges, errors.New("not implemented")
}
