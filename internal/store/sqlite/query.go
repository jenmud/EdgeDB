package sqlite

import (
	"context"
	"errors"

	"github.com/jenmud/edgedb/internal/store/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Query struct {
	db *sqlx.DB
}

// New creates a new Query instance with the provided database connection.
func New(dns string) *Query {
	return &Query{db: sqlx.MustConnect("sqlite3", dns)}
}

func (q *Query) InsertNode(ctx context.Context, node models.Node) (models.Node, error) {
	return models.Node{}, errors.New("not implemented")
}

func (q *Query) InsertEdge(ctx context.Context, edge models.Edge) (models.Edge, error) {
	return models.Edge{}, errors.New("not implemented")
}

func (q *Query) Nodes(ctx context.Context) ([]models.Node, error) {
	var nodes []models.Node
	return nodes, errors.New("not implemented")
}

func (q *Query) Edges(ctx context.Context) ([]models.Edge, error) {
	var edges []models.Edge
	return edges, errors.New("not implemented")
}
