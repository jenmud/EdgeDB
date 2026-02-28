package models

import (
	"context"
	"time"
)

// Node represents a node in the store.
type Node struct {
	ID         uint64     `db:"id" json:"id"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	Label      string     `db:"label" json:"label"`
	Properties Properties `db:"properties" json:"properties"`
	Snippet    string     `db:"-" json:"snippet,omitempty"` // this is a special field show a small snippet of the match terms
}

// NewNode creates a new node with the given label and properties.
func NewNode(ctx context.Context, label string, properties Properties) Node {
	return Node{Label: label, Properties: properties}
}
