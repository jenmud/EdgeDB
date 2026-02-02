package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jenmud/edgedb/internal/store/models"
	"github.com/jmoiron/sqlx"
)

// DB extends sqlx.DB to implement the Querier interface.
// A implementation of the Querier interface can embed this struct to inherit its methods.
type DB struct {
	*sqlx.DB
}

func (b *DB) Tx(ctx context.Context) (*sql.Tx, error) {
	return b.DB.BeginTx(ctx, nil)
}

func (b *DB) Close() error {
	return b.DB.Close()
}

func (b *DB) InsertNode(ctx context.Context, node models.Node) (models.Node, error) {
	return models.Node{}, errors.New("not implemented")
}

func (b *DB) InsertEdge(ctx context.Context, edge models.Edge) (models.Edge, error) {
	return models.Edge{}, errors.New("not implemented")
}

func (b *DB) Nodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	return nodes, errors.New("not implemented")
}

func (b *DB) Edges(ctx context.Context) ([]models.Edge, error) {
	var edges []models.Edge
	return edges, errors.New("not implemented")
}
