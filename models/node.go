package models

import (
	"context"
)

// Node represents a node in the store.
type Node struct {
	ID         uint64     `db:"id" json:"id"`
	Label      string     `db:"label" json:"label"`
	Properties Properties `db:"properties" json:"properties"`
}

// NewNode creates a new node with the given label and properties.
func NewNode(ctx context.Context, label string, properties Properties) Node {
	return Node{Label: label, Properties: properties}
}
