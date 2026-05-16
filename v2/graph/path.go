package graph

import "fmt"

// Path represents a path through the graph, consisting of a sequence of node IDs and the total distance (weight) of the path.
type Path struct {
	distances []distance
}

// distance represents a single hop in the path, containing the destination node and the weight of the edge to that node.
type distance struct {
	node   NodeID
	weight Weight
}

// Path returns a Path object representing the path through the graph defined by the given sequence of node IDs.
// If the path is valid (i.e., there are edges between all consecutive nodes), it returns a Path with the total
// distance and the sequence of nodes. If any edge in the path does not exist, it returns a Path marked as
// infinite (unreachable) and an error indicating which edge is missing.
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
