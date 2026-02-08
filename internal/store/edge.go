package store

// Direction is the edge direction
type Direction int

const (
	// Both is both incoming and outgoing edges - 0
	Both Direction = iota
	// In is only incoming edges - 1
	In
	// Out is only outgoing edges - 2
	Out
)

// Edge represents an edge in the store.
type Edge struct {
	ID         uint64     `db:"id" json:"id"`
	Label      string     `db:"label" json:"label"`
	Properties Properties `db:"properties" json:"properties"`
	FromNodes  []uint64   `db:"from_nodes" json:"from_nodes"`
	ToNodes    []uint64   `db:"to_nodes" json:"to_nodes"`
}
