package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/jenmud/edgedb/models"
	"github.com/jenmud/edgedb/pkg/common"
	"modernc.org/sqlite"
	_ "modernc.org/sqlite"

	"github.com/golang-migrate/migrate/v4"
	migrateSQLite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed "migrations/*.sql"
var migrations embed.FS
var once sync.Once

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
	once.Do(registerFuncs)

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

// registerFuncs registers custom SQL functions for SQLite. It will panic if the registration fails, so it should be called during initialization.
func registerFuncs() {
	slog.Debug("registering custom sql functions")

	sqlite.MustRegisterDeterministicScalarFunction(
		"json_extract_keys",
		1,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			var payload json.RawMessage

			switch argTyped := args[0].(type) {
			case string:
				payload = json.RawMessage([]byte(argTyped))
			case []byte:
				payload = json.RawMessage(argTyped)
			default:
				return nil, fmt.Errorf("expected argument to be a string, got: %T", argTyped)
			}

			props := make(map[string]any)
			if err := json.Unmarshal(payload, &props); err != nil {
				return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
			}

			keys := common.Keys(props)
			return strings.Join(keys, ","), nil
		},
	)

	sqlite.MustRegisterDeterministicScalarFunction(
		"json_extract_values",
		1,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			var payload json.RawMessage

			switch argTyped := args[0].(type) {
			case string:
				payload = json.RawMessage([]byte(argTyped))
			case []byte:
				payload = json.RawMessage(argTyped)
			default:
				return nil, fmt.Errorf("expected argument to be a string, got: %T", argTyped)
			}

			props := make(map[string]any)
			if err := json.Unmarshal(payload, &props); err != nil {
				return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
			}

			values := common.Values(props)
			return strings.Join(values, ","), nil
		},
	)
}

// UpsertNodes inserts or creates one or more nodes.
func UpsertNodes(ctx context.Context, tx *sql.Tx, n ...models.Node) ([]models.Node, error) {
	nodes := make([]models.Node, len(n))

	for i, n := range n {

		node := models.Node{}

		props, err := n.Properties.ToBytes()
		if err != nil {
			return nodes, err
		}

		var row *sql.Row

		if n.ID == 0 {
			query := `
				INSERT INTO nodes (label, properties)
				VALUES (?, ?)
				RETURNING id, label, properties;
			`

			row = tx.QueryRowContext(ctx, query, n.Label, props)
		} else {
			query := `
				INSERT OR REPLACE INTO nodes (id, label, properties)
				VALUES (?, ?, ?)
				RETURNING id, label, properties;
			`

			row = tx.QueryRowContext(ctx, query, n.ID, n.Label, props)
		}

		if err := row.Scan(&node.ID, &node.Label, &props); err != nil {
			return nodes, err
		}

		if err := node.Properties.FromBytes(props); err != nil {
			return nodes, err
		}

		nodes[i] = node
	}

	return nodes, nil
}
