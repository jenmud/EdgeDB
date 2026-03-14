package models

// Graph returns the graph data in format that can be used with D3/Force-graph JS.
type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// AddNodes adds one or more nodes as graph node to support the special formatting.
func (g *Graph) AddNodes(nodes ...Node) {
	for _, n := range nodes {
		g.Nodes = append(g.Nodes, n)
	}
}

// AddEdges adds one or more links/edges as graph links to support the special formatting.
func (g *Graph) AddEdges(edges ...Edge) {
	for _, e := range edges {
		if e.From == 0 || e.To == 0 {
			continue
		}
		g.Edges = append(g.Edges, e)
	}
}
