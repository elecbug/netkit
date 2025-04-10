package graph

// `Path` represents a path between two nodes in a graph.
// It contains the total distance (sum of edge weights) for the path
// and a sequence of nodes (Identifiers) traversed in the path.
type Path struct {
	distance Distance // Total distance or weight of the path.
	nodes    []NodeID // Sequence of nodes from the source to the destination.
}

// `NewPath` creates and initializes a new Path instance.
//
// Parameters:
//   - distance: The total weight or distance of the path.
//   - nodes: A slice of node identifiers representing the sequence of nodes in the path.
// Returns a pointer to the newly created Path.
func NewPath(distance Distance, nodes []NodeID) *Path {
	return &Path{
		distance: distance,
		nodes:    nodes,
	}
}

// `Distance` returns the total distance of the path.
// This is the sum of all edge weights along the path.
func (p Path) Distance() Distance {
	return p.distance
}

// `Nodes` returns the sequence of nodes in the path.
// The slice represents the order of traversal from the source node to the destination node.
func (p Path) Nodes() []NodeID {
	return p.nodes
}
