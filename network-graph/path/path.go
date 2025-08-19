// Package path defines path structures used by network-graph algorithms.
package path

import (
	"github.com/elecbug/netkit/network-graph/node"
)

// Path represents an ordered sequence of nodes with a hop distance.
type Path struct {
	distance int
	nodes    []node.ID
	isInf    bool
}

// GraphPaths is a mapping of start node IDs to end node IDs and their corresponding paths.
type GraphPaths map[node.ID]map[node.ID][]Path

// PathLength represents the length of a path between two nodes.
type PathLength map[node.ID]map[node.ID]int

// New constructs a Path from the given nodes. Distance is hops (edges).
// If no nodes are provided, the path is considered infinite (unreachable).
func New(nodes ...node.ID) *Path {
	return &Path{
		distance: len(nodes) - 1, // Assuming distance is the number of edges
		isInf:    len(nodes) == 0,
		nodes:    nodes,
	}
}

func NewSelf(id node.ID) *Path {
	return &Path{
		distance: 0,
		isInf:    false,
		nodes:    []node.ID{id},
	}
}

// IsInfinite reports whether the path is infinite (unreachable).
func (p *Path) IsInfinite() bool {
	return p.isInf
}

// Distance returns the hop distance of the path.
func (p *Path) Distance() int {
	return p.distance
}

// Nodes returns the node IDs in the path.
func (p *Path) Nodes() []node.ID {
	return p.nodes
}

// OnlyLength returns a slice of PathLength representing the lengths of all paths in the graph.
func (g GraphPaths) OnlyLength() PathLength {
	results := make(PathLength, 0)

	for start, endMap := range g {
		for end, paths := range endMap {
			if len(paths) == 0 {
				continue
			}

			if results[start] == nil {
				results[start] = make(map[node.ID]int)
			}

			results[start][end] = paths[0].Distance()
		}
	}

	return results
}
