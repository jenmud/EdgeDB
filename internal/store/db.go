package store

import (
	"context"

	"github.com/jenmud/edgedb/internal/store/models"
)

type Querier interface {
	// ApplyMigrations applies any database migrations.
	ApplyMigrations(ctx context.Context) error

	// InsertNode inserts a new node into the database.
	InsertNode(context.Context, models.Node) (models.Node, error)

	// InsertEdge inserts a new edge into the database.
	InsertEdge(context.Context, models.Edge) (models.Edge, error)

	// Nodes retrieves all nodes from the database.
	Nodes(context.Context) ([]models.Node, error)

	// Edges retrieves all edges from the database.
	Edges(context.Context) ([]models.Edge, error)

	// Node retrieves a single node by its ID.
	Node(context.Context, uint64) (models.Node, error)

	// Edge retrieves a single edge by its ID.
	Edge(context.Context, uint64) (models.Edge, error)
}
