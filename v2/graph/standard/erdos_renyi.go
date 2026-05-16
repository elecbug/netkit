package standard

import (
	"fmt"

	"github.com/elecbug/netkit/v2/graph"
)

// ErdosRenyiGraph generates a random graph based on the Erdős-Rényi model.
// Nodes are added first, and then edges are created between nodes with a probability p.
func ErdosRenyiGraph(seed int, directed bool, weightFunc WeightedFunc, n int, p float64) (*graph.Graph, error) {
	if p < 0 || p > 1 {
		return nil, fmt.Errorf("invalid probability: p must be between 0 and 1")
	}
	if n < 0 {
		return nil, fmt.Errorf("invalid number of nodes: n must be non-negative")
	}
	if weightFunc == nil {
		weightFunc = Unweighted()
	}

	r := generateRand(seed)
	g := graph.New(directed, true)

	for i := 0; i < n; i++ {
		if err := g.AddNode(graph.NodeID(fmt.Sprintf("%d", i))); err != nil {
			return nil, fmt.Errorf("failed to add node: %w", err)
		}
	}

	if !directed {
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if r.Float64() < p {
					from := graph.NodeID(fmt.Sprintf("%d", i))
					to := graph.NodeID(fmt.Sprintf("%d", j))
					if err := g.AddEdge(from, to, weightFunc(from, to)); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			}
		}
	} else {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j && r.Float64() < p {
					from := graph.NodeID(fmt.Sprintf("%d", i))
					to := graph.NodeID(fmt.Sprintf("%d", j))
					if err := g.AddEdge(from, to, weightFunc(from, to)); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			}
		}
	}

	return g, nil
}
