package store

import (
	"context"
	"encoding/json"
)

// Node represents a node in the store.
type Node struct {
	ID         uint64         `db:"id" json:"id"`
	Name       string         `db:"name" json:"name"`
	Properties map[string]any `db:"properties" json:"properties"`
}

// Sync updates or inserts the Node in the provided store.
func (n Node) Sync(ctx context.Context, s *DB) error {
	tx, err := s.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var props json.RawMessage

	if n.Properties != nil {
		b, err := json.Marshal(n.Properties)
		if err != nil {
			return err
		}

		props = b
	}

	query := `
		INSERT INTO nodes (ID, name, properties)
		VALUES (?, ?, ?)
		ON CONFLICT (id)
		DO UPDATE SET
			name = EXCLUDED.name,
			properties = EXCLUDED.properties;
	`

	row := tx.QueryRowContext(ctx, query, n.ID, n.Name, props)

	if err := row.Scan(&n.ID, &n.Name, &props); err != nil {
		return err
	}

	if err := json.Unmarshal(props, &n.Properties); err != nil {
		return err
	}

	return tx.Commit()
}
