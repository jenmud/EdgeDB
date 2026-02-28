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
	"time"

	"github.com/jenmud/edgedb/internal/store"
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

// New creates a new store instance with the provided database connection.
func New(ctx context.Context, dns string) (*Store, error) {
	s := &Store{}

	db, err := sql.Open("sqlite", dns)
	if err != nil {
		return nil, err
	}

	s.db = db

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

	return s, ApplyMigrations(ctx, s.db)
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

// Store is the underlying sqlite store.
type Store struct {
	db *sql.DB
}

// Close closed the store.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Tx returns a new transaction.
func (s *Store) Tx(ctx context.Context) (*sql.Tx, error) {
	if s.db == nil {
		return nil, errors.New("no attached database found")
	}
	return s.db.BeginTx(ctx, nil)
}

// UpsertNodes inserts or creates one or more nodes.
func (s *Store) UpsertNodes(ctx context.Context, n ...models.Node) ([]models.Node, error) {
	tx, err := s.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	nodes := make([]models.Node, len(n))

	for i, n := range n {

		node := models.Node{}

		props, err := n.Properties.ToBytes()
		if err != nil {
			return nodes, err
		}

		// We need to pass in a null ID id the node ID 0
		// so that the database can assign a new ID.
		var id *uint64

		if n.ID <= 0 {
			// DB will assign a new ID
			id = nil
		} else {
			// DB will either insert with this ID or update an existing Node if the ID conflicts
			id = &n.ID
		}

		query := `
			INSERT INTO nodes (id, label, properties)
			VALUES (?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				id = excluded.id,
				label = excluded.label,
				properties = excluded.properties
			RETURNING id, label, properties;
		`

		row := tx.QueryRowContext(ctx, query, id, n.Label, props)

		if err := row.Scan(&node.ID, &node.Label, &props); err != nil {
			return nodes, err
		}

		if err := node.Properties.FromBytes(props); err != nil {
			return nodes, err
		}

		nodes[i] = node
	}

	return nodes, tx.Commit()
}

// DefaultLimit is the default limit of return items to return.
const DefaultLimit int = 1000

// NodesTermSearch applies the search term and returns nodes with match. Limit defaults to 1000 if limit is 0
func (s *Store) NodesTermSearch(ctx context.Context, args store.NodesTermSearchArgs) ([]models.Node, error) {
	if args.Limit == 0 {
		args.Limit = DefaultLimit
	}

	if args.SnippetTokens < 0 {
		args.SnippetTokens = 10
	}

	if args.SnippetTokens > 64 {
		args.SnippetTokens = 64
	}

	if args.SnippetStart == "" {
		args.SnippetStart = `<span class="text-red-500">`
	}

	if args.SnippetEnd == "" {
		args.SnippetEnd = `</span>`
	}

	query := `
	SELECT n.id, n.created_at, n.updated_at, n.label, n.properties, snippet(fts, -1, ?, ?, ' ... ', ?) as snippet
	FROM fts
	JOIN nodes n ON n.id = fts.id
	WHERE fts.type = 'node' AND fts MATCH ?
	ORDER BY bm25(fts)
	LIMIT ?;
	`

	rows, err := s.db.QueryContext(ctx, query, args.SnippetStart, args.SnippetEnd, args.SnippetTokens, args.Term, args.Limit)
	if err != nil {
		return nil, err
	}

	nodes := []models.Node{}

	for rows.Next() {
		n := models.Node{}

		var createdAt int64
		var updatedAt int64

		var props []byte
		if err := rows.Scan(&n.ID, &createdAt, &updatedAt, &n.Label, &props, &n.Snippet); err != nil {
			return nodes, err
		}

		if err := n.Properties.FromBytes(props); err != nil {
			return nodes, err
		}

		n.CreatedAt = time.Unix(createdAt, 0)
		n.UpdatedAt = time.Unix(updatedAt, 0)

		nodes = append(nodes, n)
	}

	return nodes, nil
}

// Nodes applies the search for all nodes in the store.
func (s *Store) Nodes(ctx context.Context, args store.NodesArgs) ([]models.Node, error) {
	if args.Limit == 0 {
		args.Limit = DefaultLimit
	}

	query := `
	SELECT n.id, n.created_at, n.updated_at, n.label, n.properties
	FROM nodes n
	LIMIT ?;
	`

	rows, err := s.db.QueryContext(ctx, query, args.Limit)
	if err != nil {
		return nil, err
	}

	nodes := []models.Node{}

	for rows.Next() {
		n := models.Node{}

		var createdAt int64
		var updatedAt int64

		var props []byte
		if err := rows.Scan(&n.ID, &createdAt, &updatedAt, &n.Label, &props); err != nil {
			return nodes, err
		}

		if err := n.Properties.FromBytes(props); err != nil {
			return nodes, err
		}

		n.CreatedAt = time.Unix(createdAt, 0)
		n.UpdatedAt = time.Unix(updatedAt, 0)

		nodes = append(nodes, n)
	}

	return nodes, nil
}
