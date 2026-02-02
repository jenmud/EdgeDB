package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jenmud/edgedb/internal/store/models"
	"github.com/jmoiron/sqlx"
)

// DB extends sqlx.DB implementing additional methods used for the store..
type DB struct {
	*sqlx.DB
}

// Tx starts a new transaction.
func (b *DB) Tx(ctx context.Context) (*sql.Tx, error) {
	return b.DB.BeginTx(ctx, nil)
}

// Close closed the database.
func (b *DB) Close() error {
	return b.DB.Close()
}

// InsertNode inserts a new node into the store.
func (b *DB) InsertNode(ctx context.Context, node models.Node) (models.Node, error) {
	return models.Node{}, errors.New("not implemented")
}

// InsertEdite inserts a new edge into the store.
func (b *DB) InsertEdge(ctx context.Context, edge models.Edge) (models.Edge, error) {
	return models.Edge{}, errors.New("not implemented")
}

// Nodes returns all the nodes in the store.
func (b *DB) Nodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	return nodes, errors.New("not implemented")
}

// Edites returns all the edges in the store.
func (b *DB) Edges(ctx context.Context) ([]models.Edge, error) {
	var edges []models.Edge
	return edges, errors.New("not implemented")
}
