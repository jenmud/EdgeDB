package models

import (
	"context"
	"time"
)

type Edge struct {
	ID         uint64     `db:"id" json:"id"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	Label      string     `db:"label" json:"label"`
	Properties Properties `db:"properties" json:"properties"`
	From       uint64     `db:"from_id" json:"from_id"`
	To         uint64     `db:"to_id" json:"to_id"`
	Weight     int        `db:"weight" json:"weight"`
	Snippet    string     `db:"-" json:"snippet,omitempty"` // this is a special field show a small snippet of the match terms
}

// NewEdge returns a new edge linking two nodes together.
func NewEdge(ctx context.Context, from uint64, label string, to uint64, weight int, properties Properties) Edge {
	return Edge{Label: label, From: from, To: to, Weight: weight, Properties: properties}
}
