package store

import (
	"context"
	"errors"
)

// Node represents a node in the store.
type Node struct {
	db         *DB
	ID         uint64     `db:"id" json:"id"`
	Label      string     `db:"label" json:"label"`
	Properties Properties `db:"properties" json:"properties"`
}

// NewNode creates a new node with the given label and properties.
// If you attache a store to the node using WithStore, you can call `.Sync(ctx)` to persist the node to the store.
func NewNode(ctx context.Context, store *DB, label string, properties Properties) (*Node, error) {
	n := &Node{
		Label:      label,
		Properties: properties,
	}

	if store != nil {
		n.db = store
		if err := n.Sync(ctx); err != nil {
			return nil, err
		}
		return n, nil
	}

	return n, nil
}

// WithStore sets the internal store for the node returning the node with the store attached.
func (n *Node) WithStore(store *DB) *Node {
	n.db = store
	return n
}

// Sync will sync the node with the store if it is bound to a store.
// It wll insert the node if is it new, or replace/update the node if it already exists.
func (n *Node) Sync(ctx context.Context) error {
	if n.db == nil {
		return errors.New("node is not bound to a store")
	}

	updated, err := n.db.SyncNodes(ctx, n)
	if err != nil {
		return err
	}

	if len(updated) != 1 {
		return errors.New("unexpected number of nodes returned from SyncNodes")
	}

	n.ID = updated[0].ID
	n.Label = updated[0].Label
	n.Properties = updated[0].Properties

	return nil
}
