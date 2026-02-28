package models

import "time"

type Edge struct {
	ID         uint64     `db:"id" json:"id"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	Label      string     `db:"label" json:"label"`
	Properties Properties `db:"properties" json:"properties"`
	FromID     uint64     `db:"from_id" json:"from_id"`
	ToID       uint64     `db:"to_id" json:"to_id"`
	Weight     int64      `db:"weight" json:"weight"`
	Snippet    string     `db:"-" json:"snippet,omitempty"` // this is a special field show a small snippet of the match terms
}
