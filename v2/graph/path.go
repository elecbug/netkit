package graph

import "fmt"

// Path represents a sequence of nodes in the graph along with the total distance (weight) of the path.
// It contains a slice of distances, where each distance represents a hop from one node to the next, including the weight of that hop.
type Path struct {
	distances []distance
}

// distance represents a single hop in the path, containing the destination node and the weight of the edge to that node.
type distance struct {
	node   NodeID
	weight Weight
}

// Path builds a Path for the given node sequence.
// It returns an error if any consecutive edge does not exist.
func (g *Graph) Path(nodes ...NodeID) (*Path, error) {
	if len(nodes) == 0 {
		return &Path{
			distances: []distance{},
		}, nil
	}

	path := &Path{
		distances: make([]distance, 0, len(nodes)),
	}
	path.distances = append(path.distances, distance{
		node:   nodes[0],
		weight: 0,
	})

	for i := 0; i < len(nodes)-1; i++ {
		from := nodes[i]
		to := nodes[i+1]

		weight, err := g.EdgeWeight(from, to)
		if err != nil {
			return &Path{
				distances: []distance{},
			}, fmt.Errorf("no edge from %s to %s: %w", from, to, err)
		}

		path.distances = append(path.distances, distance{
			node:   to,
			weight: weight,
		})
	}

	return path, nil
}

// Nodes returns the sequence of node IDs in the path, including the starting node.
func (p *Path) Nodes() []NodeID {
	if len(p.distances) == 0 {
		return []NodeID{}
	}

	nodes := make([]NodeID, 0, len(p.distances))

	for _, d := range p.distances {
		nodes = append(nodes, d.node)
	}

	return nodes
}

// TotalDistance returns the total distance (weight) of the path.
func (p *Path) TotalDistance() Weight {
	total := Weight(0)
	for _, d := range p.distances {
		total += d.weight
	}
	return total
}
