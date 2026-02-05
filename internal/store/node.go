package store

import (
	"context"
	"errors"
)

// Node represents a node in the store.
type Node struct {
	db         *DB
	ID         uint64     `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	Properties Properties `db:"properties" json:"properties"`
}

// Bind binds the node to the provided store.
func (n *Node) Bind(store *DB) error {
	return errors.New("not implemented")
}

// Sync will sync the node with the store if it is bound to a store.
// It wll insert the node if is it new, or replace/update the node if it already exists.
func (n *Node) Sync(ctx context.Context) error {
	return errors.New("not implemented")
}

// Edges will return all the attached edges.
func (n *Node) Edges(ctx context.Context, direction Direction, limit uint) ([]Edge, error) {
	return nil, errors.New("not implemented")
}
