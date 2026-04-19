package models

// Graph returns the graph data in format that can be used with D3/Force-graph JS.
type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}
