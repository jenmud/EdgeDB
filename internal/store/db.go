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
	db *sqlx.DB
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
		return &DB{db: db}, sqlite.ApplyMigrations(ctx, db.DB)
	}

	return nil, errors.New("unsupported store")
}

// Close closed the store.
func (b *DB) Close() error {
	return b.db.Close()
}

// Tx returns a new transaction. You must `.Commit` or `.Rollback` when you are done with the transaction.
func (b *DB) Tx(ctx context.Context) (*sql.Tx, error) {
	return b.db.BeginTx(ctx, nil)
}

// InsertNode inserts a new node into the store.
func (b *DB) InsertNode(ctx context.Context, name string, props Properties) (Node, error) {
	nodes, err := b.SyncNodes(ctx, Node{Label: name, Properties: props})
	if err != nil {
		return Node{}, err
	}

	if len(nodes) != 1 {
		return Node{}, sql.ErrNoRows
	}

	return nodes[0], nil
}

// SyncNodes syncs one or more nodes with the store.
// The node will be create in the store, but if conflict it will do a replace.
func (b *DB) SyncNodes(ctx context.Context, nodes ...Node) ([]Node, error) {
	inserted := make([]Node, 0, len(nodes))

	tx, err := b.Tx(ctx)
	if err != nil {
		return inserted, err
	}

	defer tx.Rollback()

	for _, n := range nodes {
		f := insertNode

		if n.ID > 0 {
			f = upsertNode
		}

		node, err := f(ctx, tx, n)
		if err != nil {
			return inserted, err
		}

		inserted = append(inserted, node)
	}

	return inserted, tx.Commit()
}

// insertNode inserts a new node to the store.
func insertNode(ctx context.Context, tx *sql.Tx, n Node) (Node, error) {
	var node Node

	props, err := n.Properties.ToBytes()
	if err != nil {
		return node, err
	}

	// TODO: this statement should come from the driver used.
	query := `
		INSERT INTO nodes (label, properties)
		VALUES (?, ?)
		RETURNING id, label, properties;
	`

	row := tx.QueryRowContext(ctx, query, n.Label, props)

	if err := row.Scan(&node.ID, &node.Label, &props); err != nil {
		return node, err
	}

	if err := node.Properties.FromBytes(props); err != nil {
		return node, err
	}

	return node, err
}

// upsertNode inserts or create a node in the store using the provided ID attached to the node.
func upsertNode(ctx context.Context, tx *sql.Tx, n Node) (Node, error) {
	var node Node

	props, err := n.Properties.ToBytes()
	if err != nil {
		return node, err
	}

	// TODO: this statement should come from the driver used.
	query := `
		INSERT OR REPLACE INTO nodes (id, label, properties)
		VALUES (?, ?, ?)
		RETURNING id, label, properties;
	`

	row := tx.QueryRowContext(ctx, query, n.ID, n.Label, props)

	if err := row.Scan(&node.ID, &node.Label, &props); err != nil {
		return node, err
	}

	if err := node.Properties.FromBytes(props); err != nil {
		return node, err
	}

	return node, err
}

// Nodes returns all the nodes in the store.
func (b *DB) NodeByID(ctx context.Context, id uint64) (Node, error) {
	query := `
		SELECT * FROM nodes
		WHERE id = ?;
	`

	var node Node
	err := b.db.GetContext(ctx, &node, query, id)

	return node, err
}

// safetyLimit is a safeguard limit
const safetyLimit = 1000

// validateLimit is a function to validate and return a safe limit for returning multiple search items.
func validateLimit(limit uint) uint {
	if limit == 0 || limit > safetyLimit {
		return safetyLimit
	}

	return limit
}

// Nodes returns all the nodes in the store.
// limit defaults to `safetyLimit` (see const above) if ==0 or >safetyLimit.
func (b *DB) Nodes(ctx context.Context, limit uint) ([]Node, error) {

	limit = validateLimit(limit)
	nodes := make([]Node, 0, limit)

	query := `
		SELECT * FROM nodes
		LIMIT ?;
	`

	return nodes, b.db.SelectContext(ctx, &nodes, query, limit)
}

// Edges returns all the edges in the store.
func (b *DB) Edges(ctx context.Context) ([]Edge, error) {
	var edges []Edge
	return edges, errors.New("not implemented")
}
