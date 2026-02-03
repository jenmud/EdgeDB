package store

import (
	"context"
	"database/sql"
	"encoding/json"
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

// Close closed the store.
func (b *DB) Close() error {
	return b.DB.Close()
}

// InsertNode inserts a new node into the store.
func (b *DB) InsertNode(ctx context.Context, n Node) (Node, error) {
	nodes, err := b.InsertNodes(ctx, n)
	if err != nil {
		return Node{}, err
	}

	if len(nodes) != 1 {
		return Node{}, sql.ErrNoRows
	}

	return nodes[0], nil
}

// InsertNodes inserts one or more nodes into the store.
func (b *DB) InsertNodes(ctx context.Context, nodes ...Node) ([]Node, error) {
	inserted := make([]Node, 0, len(nodes))

	tx, err := b.BeginTxx(ctx, nil)
	if err != nil {
		return inserted, err
	}

	defer tx.Rollback()

	for _, n := range nodes {

		var node Node
		var props json.RawMessage

		if n.Properties != nil {
			propsBytes, err := json.Marshal(n.Properties)
			if err != nil {
				return inserted, err
			}

			props = propsBytes
		}

		// TODO: this statement should come from the driver used.
		query := `
			INSERT INTO nodes (name, properties)
			VALUES (?, ?)
			RETURNING id, name, properties;
		`

		row := tx.QueryRowContext(ctx, query, n.Name, props)

		if err := row.Scan(&node.ID, &node.Name, &props); err != nil {
			return inserted, err
		}

		if err := json.Unmarshal(props, &node.Properties); err != nil {
			return inserted, err
		}

		inserted = append(inserted, node)
	}

	return inserted, tx.Commit()
}

// Nodes returns all the nodes in the store.
func (b *DB) NodeByID(ctx context.Context, id uint64) (Node, error) {
	query := `
		SELECT * FROM nodes
		WHERE id = ?;
	`

	var node Node
	err := b.GetContext(ctx, &node, query, id)

	return node, err
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
