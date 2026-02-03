package store

// Edge represents an edge in the store.
type Edge struct {
	ID         uint64         `db:"id" json:"id"`
	Name       string         `db:"name" json:"name"`
	Properties map[string]any `db:"properties" json:"properties"`
	FromNodes  []uint64       `db:"from_nodes" json:"from_nodes"`
	ToNodes    []uint64       `db:"to_nodes" json:"to_nodes"`
}
