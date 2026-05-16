package standard

import (
	"fmt"

	"github.com/elecbug/netkit/v2/graph"
)

// BarabasiAlbertGraph generates a graph based on the Barabási-Albert preferential attachment model.
// The graph starts with m fully connected nodes, and new nodes are added one by one, each connecting
// to m existing nodes with a probability proportional to their degree.
func BarabasiAlbertGraph(seed int, directed bool, weightFunc WeightedFunc, n int, m int) (*graph.Graph, error) {
	if n < 0 {
		return nil, fmt.Errorf("invalid number of nodes: n must be non-negative")
	}
	if m < 1 || n <= m {
		return nil, fmt.Errorf("invalid parameters: n must be greater than m and m must be at least 1")
	}
	if weightFunc == nil {
		weightFunc = Unweighted()
	}

	r := generateRand(seed)
	g := graph.New(directed, true)

	// --- 1. initialize ---
	for i := 0; i < m; i++ {
		if err := g.AddNode(graph.NodeID(fmt.Sprintf("%d", i))); err != nil {
			return nil, fmt.Errorf("failed to add node: %w", err)
		}
	}
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			from := graph.NodeID(fmt.Sprintf("%d", i))
			to := graph.NodeID(fmt.Sprintf("%d", j))
			if err := g.AddEdge(from, to, weightFunc(from, to)); err != nil {
				return nil, fmt.Errorf("failed to add edge: %w", err)
			}
		}
	}

	// --- 2. preferential attachment ---
	for i := m; i < n; i++ {
		newNode := graph.NodeID(fmt.Sprintf("%d", i))
		if err := g.AddNode(newNode); err != nil {
			return nil, fmt.Errorf("failed to add node: %w", err)
		}

		// calculate current node degrees
		degrees := make(map[graph.NodeID]int)
		totalDegree := 0
		for _, id := range g.Nodes() {
			node, err := g.Node(id)
			if err != nil {
				return nil, fmt.Errorf("failed to get node: %w", err)
			}
			d := node.Degree() // each node degree
			degrees[id] = d
			totalDegree += d
		}

		// degree based sampling
		chosen := make(map[graph.NodeID]bool)
		for len(chosen) < m {
			r := r.Intn(totalDegree)
			accum := 0
			var target graph.NodeID
			for id, d := range degrees {
				accum += d
				if r < accum {
					target = id
					break
				}
			}
			// self-loop and duplicate edges are not allowed
			if target != newNode && !chosen[target] {
				if err := g.AddEdge(newNode, target, weightFunc(newNode, target)); err != nil {
					return nil, fmt.Errorf("failed to add edge: %w", err)
				}
				chosen[target] = true
			}
		}
	}

	return g, nil
}
