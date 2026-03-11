package models

import "fmt"

// Graph returns the graph data in format that can be used with D3/Force-graph JS.
type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Links []GraphEdge `json:"links"`
}

// AddNodes adds one or more nodes as graph node to support the special formatting.
func (g *Graph) AddNodes(nodes ...Node) {
	for _, n := range nodes {
		g.Nodes = append(g.Nodes, GraphNode{Node: n, IDstr: fmt.Sprintf("%d", n.ID)})
	}
}

// AddLinks adds one or more links/edges as graph links to support the special formatting.
func (g *Graph) AddLinks(edges ...Edge) {
	for _, e := range edges {
		if e.From == 0 || e.To == 0 {
			continue
		}
		g.Links = append(
			g.Links,
			GraphEdge{
				Edge:    e,
				IDstr:   fmt.Sprintf("%d", e.ID),
				Fromstr: fmt.Sprintf("%d", e.From),
				Tostr:   fmt.Sprintf("%d", e.To),
			},
		)
	}
}

type GraphNode struct {
	Node
	IDstr string `db:"id" json:"id"`
}

type GraphEdge struct {
	Edge
	IDstr   string `json:"id"`
	Fromstr string `json:"source"`
	Tostr   string `json:"target"`
}
