// Package path defines path structures used by network-graph algorithms.
package path

import (
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// Path represents an ordered sequence of nodes with a hop distance.
type Path struct {
	distance int
	nodes    []node.ID
	isInf    bool
}

// GraphPaths is a mapping of start node IDs to end node IDs and their corresponding paths.
type GraphPaths map[node.ID]map[node.ID][]Path

// NewPath constructs a Path from the given nodes. Distance is hops (edges).
// If no nodes are provided, the path is considered infinite (unreachable).
func NewPath(nodes ...node.ID) *Path {
	return &Path{
		distance: len(nodes) - 1, // Assuming distance is the number of edges
		isInf:    len(nodes) == 0,
		nodes:    nodes,
	}
}

// IsInfinite reports whether the path is infinite (unreachable).
func (p *Path) IsInfinite() bool {
	return p.isInf
}

// GetDistance returns the hop distance of the path.
func (p *Path) GetDistance() int {
	return p.distance
}

// GetNodes returns the node IDs in the path.
func (p *Path) GetNodes() []node.ID {
	return p.nodes
}
