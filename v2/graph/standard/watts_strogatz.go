package standard

import (
	"fmt"

	"github.com/elecbug/netkit/v2/graph"
)

// WattsStrogatzGraph generates a small-world graph based on the Watts-Strogatz model.
// If weightFunc is nil, all edges will have no weight (unweighted graph). Otherwise, weightFunc will be called for
// each new edge with the new node and the target node as arguments.
func WattsStrogatzGraph(seed int, directed bool, weightFunc WeightedFunc, n, k int, beta float64) (*graph.Graph, error) {
	if k < 0 || k >= n {
		// degree must be between 0 and n-1
		return nil, fmt.Errorf("invalid degree: k must be between 0 and n-1")
	}
	if n < 0 {
		return nil, fmt.Errorf("invalid number of nodes: n must be non-negative")
	}
	if beta < 0 || beta > 1 {
		return nil, fmt.Errorf("invalid rewiring probability: beta must be between 0 and 1")
	}

	r := generateRand(seed)
	g := graph.New(directed, weightFunc != nil)

	if weightFunc == nil {
		weightFunc = func(from, to *graph.Node) *graph.Weight {
			return nil
		}
	}

	// --- 1. Generate Nodes ---
	for i := 0; i < n; i++ {
		if err := g.AddNode(graph.NodeID(fmt.Sprintf("%d", i))); err != nil {
			return nil, fmt.Errorf("failed to add node: %w", err)
		}
	}

	// --- 2. Generate Regular Ring Lattice ---
	for i := 0; i < n; i++ {
		for j := 1; j <= k/2; j++ {
			neighbor := (i + j) % n
			from := graph.NodeID(fmt.Sprintf("%d", i))
			to := graph.NodeID(fmt.Sprintf("%d", neighbor))
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

	// --- 3. Rewiring Phase ---
	for i := 0; i < n; i++ {
		for j := 1; j <= k/2; j++ {
			neighbor := (i + j) % n
			if r.Float64() < beta {
				// Remove existing edge
				from := graph.NodeID(fmt.Sprintf("%d", i))
				to := graph.NodeID(fmt.Sprintf("%d", neighbor))
				if err := g.RemoveEdge(from, to); err != nil {
					return nil, fmt.Errorf("failed to remove edge: %w", err)
				}

				// Select a random other node (self-loop, duplicate prevention)
				for {
					newNeighbor := graph.NodeID(fmt.Sprintf("%d", r.Intn(n)))
					if newNeighbor != from && !g.HasEdge(from, newNeighbor) {
						fromNode, err := g.Node(from)
						if err != nil {
							return nil, fmt.Errorf("failed to get node: %w", err)
						}
						newNeighborNode, err := g.Node(newNeighbor)
						if err != nil {
							return nil, fmt.Errorf("failed to get node: %w", err)
						}
						if err := g.AddEdge(from, newNeighbor, weightFunc(fromNode, newNeighborNode)); err != nil {
							return nil, fmt.Errorf("failed to add edge: %w", err)
						}
						break
					}
				}
			}
		}
	}

	return g, nil
}
