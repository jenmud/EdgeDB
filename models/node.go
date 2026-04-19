package models

import (
	"context"
	"time"
)

// Node represents a node in the store.
type Node struct {
	ID         uint64      `db:"id" json:"id"`
	CreatedAt  time.Time   `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt  time.Time   `db:"updated_at,omitempty" json:"updated_at"`
	Labels     []string    `db:"labels" json:"labels"`
	Properties map[any]any `db:"properties,omitempty" json:"properties,omitempty"`
	Snippet    string      `db:"-" json:"snippet,omitempty"` // this is a special field show a small snippet of the match terms
}
