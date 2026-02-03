package store

// Node represents a node in the store.
type Node struct {
	ID         uint64         `db:"id" json:"id"`
	Name       string         `db:"name" json:"name"`
	Properties map[string]any `db:"properties" json:"properties"`
}
