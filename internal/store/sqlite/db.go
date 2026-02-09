package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"log/slog"
	"strings"

	"github.com/jenmud/edgedb/models"
	"github.com/jenmud/edgedb/pkg/common"
	_ "modernc.org/sqlite"

	"github.com/golang-migrate/migrate/v4"
	migrateSQLite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed "migrations/*.sql"
var migrations embed.FS

// New creates a new Query instance with the provided database connection.
func New(ctx context.Context, dns string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dns)
	if err != nil {
		return nil, err
	}

	// call SetMaxOpenConns to 1 for SQLite to avoid "database is locked" errors on the original underlying DB
	db.SetMaxOpenConns(1)

	slog.SetDefault(
		slog.With(
			slog.Group(
				"store",
				slog.String("driver", "sqlite"),
				slog.String("dsn", dns),
			),
		),
	)

	slog.Debug("attached to store")
	return db, ApplyMigrations(ctx, db)
}

// ApplyMigrations applies database migrations from the embedded filesystem.
func ApplyMigrations(ctx context.Context, db *sql.DB) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		slog.Error("error loading migrations", slog.String("reason", err.Error()))
		return err
	}

	// db.DB.DB is a bit of inheritance mess
	driver, err := migrateSQLite.WithInstance(db, &migrateSQLite.Config{})
	if err != nil {
		slog.Error("error creating migration driver", slog.String("reason", err.Error()))
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite", driver)
	if err != nil {
		slog.Error("error creating migrate instance", slog.String("reason", err.Error()))
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("error applying migrations", slog.String("reason", err.Error()))
		return err
	}

	slog.Debug("migrations successfully applied")
	return nil
}

// TODO: add in a custom function to sqlite to extract the property keys and values, then we should be able to use triggers to update the fts

// UpsertNodes inserts or creates one or more nodes.
func UpsertNodes(ctx context.Context, tx *sql.Tx, n ...models.Node) ([]models.Node, error) {
	nodes := make([]models.Node, len(n))

	for i, n := range n {

		node := models.Node{}

		props, err := n.Properties.ToBytes()
		if err != nil {
			return nil, err
		}

		query := `
			INSERT OR REPLACE INTO nodes (id, label, properties)
			VALUES (?, ?, ?)
			RETURNING id, label, properties;
		`

		row := tx.QueryRowContext(ctx, query, n.ID, n.Label, props)

		if err := row.Scan(&node.ID, &node.Label, &props); err != nil {
			return nil, err
		}

		if err := node.Properties.FromBytes(props); err != nil {
			return nil, err
		}

		fts_query := `
			DELETE from nodes_fts where id = ?;
			INSERT INTO nodes_fts(id, label, prop_keys, prop_values) VALUES(?, ?, ?, ?);
		`

		keys, values := common.FlattenMAP(node.Properties)
		if _, err := tx.ExecContext(ctx, fts_query, node.ID, node.ID, node.Label, strings.Join(keys, ","), strings.Join(values, ",")); err != nil {
			slog.Error("failed to update node FTS", "error", err)
			return nil, err
		}

		nodes[i] = node
	}

	return nodes, nil
}
