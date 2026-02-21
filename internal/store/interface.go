package store

import (
	"context"

	"github.com/jenmud/edgedb/models"
)

// NodeWriter defines the behavior required to modify the node store.
type NodeWriter interface {
	// Upsert inserts or updates one or more nodes.
	UpsertNodes(context.Context, ...models.Node) ([]models.Node, error)
}

// NodesTermSearchArgs are arguments used for search term queries.
type NodesTermSearchArgs struct {
	// Term is the search term.
	Term string

	// Limit is the max number of items to return.
	Limit int
}

// NodesArgs are the search arguments for nodes in the store.
type NodesArgs struct {
	// Limit is the max number of items to return.
	Limit int
}

// NodeSearcher defines the behavior required to search for nodes in the store..
type NodeSearcher interface {
	// Nodes performs a search for all nodes in the store.
	Nodes(context.Context, NodesArgs) ([]models.Node, error)

	// NodesTermSearch performs a full-text or term-based search over nodes.
	NodesTermSearch(context.Context, NodesTermSearchArgs) ([]models.Node, error)
}

// NodeStore defines the behavior required to persist and search nodes.
type NodeStore interface {
	NodeWriter
	NodeSearcher
	Close() error
}

// Store defines the behavior required to persist and search a store.
type Store interface {
	NodeStore
	Close() error
}
