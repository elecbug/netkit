// Package graph provides a simple adjacency map graph for network-graph.
package graph

import (
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// Graph maintains nodes and adjacency edges.
type Graph struct {
	nodes         map[node.ID]bool
	edges         map[node.ID]map[node.ID]bool
	bidirectional bool
}

// New creates and returns an empty Graph.
func New(bidirectional bool) *Graph {
	return &Graph{
		nodes:         make(map[node.ID]bool),
		edges:         make(map[node.ID]map[node.ID]bool),
		bidirectional: bidirectional,
	}
}

func (g *Graph) IsBidirectional() bool {
	return g.bidirectional
}
