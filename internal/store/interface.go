package store

import (
	"context"

	"github.com/jenmud/edgedb/models"
)

// NodeWriter defines the behavior required to modify the node store.
type NodeWriter interface {
	// UpsertNodes inserts or updates one or more nodes.
	UpsertNodes(context.Context, ...models.Node) ([]models.Node, error)
}

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

// NodesArgs are the search arguments for nodes in the store.
type NodesArgs struct {
	// Limit is the max number of items to return.
	Limit int

	// LastID is the last know primary key/ID which will be used for fast pagination.
	LastID uint64
}

// NodeSearcher defines the behavior required to search for nodes in the store..
type NodeSearcher interface {
	// Nodes performs a search for all nodes in the store.
	Nodes(context.Context, NodesArgs) ([]models.Node, error)

	// NodesTermSearch performs a full-text or term-based search over nodes.
	NodesTermSearch(context.Context, TermSearchArgs) ([]models.Node, error)
}

// NodeStore defines the behavior required to persist and search nodes.
type NodeStore interface {
	NodeWriter
	NodeSearcher
	Node(context.Context, uint64) (models.Node, error)
	Close() error
}

// EdgeWriter defines the behavior required to modify the edge store.
type EdgeWriter interface {
	// UpsertEdges inserts or updates one or more edges.
	UpsertEdges(context.Context, ...models.Edge) ([]models.Edge, error)
}

// EdgesArgs are the search arguments for edges in the store.
type EdgesArgs struct {
	// Limit is the max number of items to return.
	Limit int

	// LastID is the last know primary key/ID which will be used for fast pagination.
	LastID uint64
}

// EdgeSearcher defines the behavior required to search for edges in the store..
type EdgeSearcher interface {
	// Nodes performs a search for all nodes in the store.
	Edges(context.Context, EdgesArgs) ([]models.Edge, error)

	// EdgesTermSearch performs a full-text or term-based search over edges.
	EdgesTermSearch(context.Context, TermSearchArgs) ([]models.Edge, error)
}

// EdgeStore defines the behavior required to persist and search edges.
type EdgeStore interface {
	EdgeWriter
	EdgeSearcher
	Edge(context.Context, uint64) (models.Edge, error)
	Close() error
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
	NodeStore
	EdgeStore
	Graph(context.Context, TermSearchArgs) (models.Graph, error)
	SubGraph(context.Context, SubGraphArgs) (models.Graph, error)
	Health(context.Context) models.Health
	Close() error
}
