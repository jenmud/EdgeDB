package store

import (
	"context"

	"github.com/jenmud/edgedb/models"
)


// TermSearchArgs are arguments used for search term queries.
type TermSearchArgs struct {
	// Term is the search term.
	Term string

	// Limit is the max number of items to return.
	Limit int

	// LastID is the last know primary key/ID which will be used for fast pagination.
	LastID uint64

	// SnippetTokens is the max tokens in the returned snipped text.
	SnippetTokens int

	// SnippetStart is the opening tag.
	SnippetStart string

	// SnippetEnd is the ending tag.
	SnippetEnd string
}


// SubGraphArgs are the arguments for building a sub graph.
type SubGraphArgs struct {
	// FromNodeID is the node ID to start building the sub graph from.
	FromNodeID uint64

	// ToNodeID is the node ID to start building the sub graph from.
	ToNodeID uint64

	// EdgeID is the edge ID to start building the sub graph from.
	EdgeID uint64

	// Limit is the max number of items to return.
	Limit int

	// LastID is the last know primary key/ID which will be used for fast pagination.
	LastID uint64
}

// Store defines the behavior required to persist and search a store.
type Store interface {
	Graph(context.Context, TermSearchArgs) (models.Graph, error)
	Health(context.Context) models.Health
	Close() error
}


// UnimplementedStore implements a NOOP store.
type UnimplementedStore struct {}

func (u *UnimplementedStore) Graph(ctx context.Context, TermSearchArgs) (models.Graph, error){
	return models.Graph{}, errors.New("not implemented")
}

func (u *UnimplementedStore) Health(ctx context.Context) models.Health {
	return models.Health{
		Status "not implemented"
	}
}


func (u *UnimplementedStore) Close() error {
	return nil, errors.New("not implemented")
}
