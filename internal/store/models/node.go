package models

import (
	"context"

	"github.com/jenmud/edgedb/internal/store"
)

// Node represents a node in the store.
type Node struct {
	ID         uint64         `db:"id" json:"id"`
	Name       string         `db:"name" json:"name"`
	Properties map[string]any `db:"properties" json:"properties"`
}

func (n Node) Create(ctx context.Context, s *store.DB) error {
	tx, err := s.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `

	`

	stmt := tx.PrepareContext(ctx, query)

	return tx.Commit()
}
