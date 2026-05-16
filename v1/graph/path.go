package graph

// Path represents an ordered sequence of nodes with a hop distance.
type Path struct {
	distance int
	nodes    []NodeID
	isInf    bool
}

// Paths is a mapping of start node IDs to end node IDs and their corresponding paths.
type Paths map[NodeID]map[NodeID][]Path

// PathLength represents the length of a path between two nodes.
type PathLength map[NodeID]map[NodeID]int

// New constructs a Path from the given nodes. Distance is hops (edges).
// If no nodes are provided, the path is considered infinite (unreachable).
func NewPath(nodes ...NodeID) *Path {
	if len(nodes) == 0 {
		return &Path{
			distance: 0,
			isInf:    true,
			nodes:    []NodeID{},
		}
	} else if len(nodes) == 1 {
		return &Path{
			distance: 0,
			isInf:    false,
			nodes:    []NodeID{nodes[0]},
		}
	} else {
		return &Path{
			distance: len(nodes) - 1, // Assuming distance is the number of edges
			isInf:    len(nodes) == 0,
			nodes:    nodes,
		}

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
func (p *Path) Nodes() []NodeID {
	return p.nodes
}

// OnlyLength returns a slice of PathLength representing the lengths of all paths in the graph.
func (g Paths) OnlyLength() PathLength {
	results := make(PathLength, 0)

	for start, endMap := range g {
		for end, paths := range endMap {
			if len(paths) == 0 {
				continue
			}

			if results[start] == nil {
				results[start] = make(map[NodeID]int)
			}

			results[start][end] = paths[0].Distance()
		}
	}

	return results
}
