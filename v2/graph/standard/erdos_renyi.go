package standard

import (
	"fmt"

	"github.com/elecbug/netkit/v2/graph"
)

// ErdosRenyiGraph generates a random graph based on the Erdős-Rényi model.
// Nodes are added first, and then edges are created between nodes with a probability p.
// If weightFunc is nil, all edges will have no weight (unweighted graph). Otherwise, weightFunc will be called for
// each new edge with the new node and the target node as arguments.
func ErdosRenyiGraph(seed int, directed bool, weightFunc WeightedFunc, n int, p float64) (*graph.Graph, error) {
	if p < 0 || p > 1 {
		return nil, fmt.Errorf("invalid probability: p must be between 0 and 1")
	}
	if n < 0 {
		return nil, fmt.Errorf("invalid number of nodes: n must be non-negative")
	}

	r := generateRand(seed)
	g := graph.New(directed, weightFunc != nil)

	if weightFunc == nil {
		weightFunc = func(from, to *graph.Node) *graph.Weight {
			return nil
		}
	}

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
					fromNode, err := g.Node(from)
					if err != nil {
						return nil, fmt.Errorf("failed to get node: %w", err)
					}
					toNode, err := g.Node(to)
					if err != nil {
						return nil, fmt.Errorf("failed to get node: %w", err)
					}
					if err := g.AddEdge(from, to, weightFunc(fromNode, toNode)); err != nil {
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
					fromNode, err := g.Node(from)
					if err != nil {
						return nil, fmt.Errorf("failed to get node: %w", err)
					}
					toNode, err := g.Node(to)
					if err != nil {
						return nil, fmt.Errorf("failed to get node: %w", err)
					}
					if err := g.AddEdge(from, to, weightFunc(fromNode, toNode)); err != nil {
						return nil, fmt.Errorf("failed to add edge: %w", err)
					}
				}
			}
		}
	}

	return g, nil
}
