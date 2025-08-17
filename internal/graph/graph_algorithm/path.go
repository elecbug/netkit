// Package graph_algorithm provides graph algorithms and analysis utilities
// built on top of the graph package.
package graph_algorithm

import (
	"github.com/elecbug/go-dspkg/graph"
)

// `Path` represents a path between two nodes in a graph.
// It contains the total distance (sum of edge weights) for the path
// and a sequence of nodes (Identifiers) traversed in the path.
type Path struct {
	distance graph.Distance // Total distance or weight of the path.
	nodes    []graph.NodeID // Sequence of nodes from the source to the destination.
}

// `NewPath` creates and initializes a new `Path` instance.
func newPath(distance graph.Distance, nodes []graph.NodeID) *Path {
	return &Path{
		distance: distance,
		nodes:    nodes,
	}
}

// `Distance` returns the total distance of the path.
// This is the sum of all edge weights along the path.
func (p Path) Distance() graph.Distance {
	return p.distance
}

// `Nodes` returns the sequence of nodes in the path.
// The slice represents the order of traversal from the source node to the destination node.
func (p Path) Nodes() []graph.NodeID {
	return p.nodes
}
