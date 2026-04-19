package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jenmud/edgedb/models"
)

// New returns a new Store.
func New(ctx context.Context, dsn string) (*Store, error) {
	conn, err := pgx.Connect(ctx, dsn)

	if err != nil {
		err = fmt.Errors("error creating store: %w", err)
		return nil, err
	}

	return &Store{db: conn}, nil
}

// Store is a graph store.
type Store struct {
	db *pgx.Conn
}

// NodesArgs are the args used fetching nodes in the store.
type NodesArgs struct {}

// Nodes returns nodes in the graph.
func (s *Store) Nodes(ctx context.Context, args NodesArgs) ([]models.Node, error){
	return models.Graph{}, errors.New("not implemented")
}

// Node returns a node from the graph.
func (s *Store) Node(ctx context.Context, id uint64) (models.Node, error){
	return models.Graph{}, errors.New("not implemented")
}

// EdgesArgs are the args used fetching edges in the store.
type EdgesArgs struct {}

// Edges returns edges in the graph.
func (s *Store) Edges(ctx context.Context, args EdgesArgs) ([]models.Edge, error){
	return models.Graph{}, errors.New("not implemented")
}

// Edge returns a edge from the graph.
func (s *Store) Edge(ctx context.Context, id uint64) (models.Edge, error){
	return models.Graph{}, errors.New("not implemented")
}

// GraphArgs are the arguments used for returning a graph.
type GraphArgs {}

// Graph returns a graph.
func (s *Store) Graph(ctx context.Context, args GraphArgs) (models.Graph, error){
	return models.Graph{}, errors.New("not implemented")
}

// /Health returns the health status.
func (s *Store) Health(ctx context.Context) models.Health {
	return models.Health{
		Status "not implemented"
	}
}

// Close closes the store.
func (s *Store) Close(ctx context.Context) error {
	return db.Close(ctx)
}
