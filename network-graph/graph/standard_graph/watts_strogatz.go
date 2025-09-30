package standard_graph

import (
	"github.com/elecbug/netkit/network-graph/graph"
	"github.com/elecbug/netkit/network-graph/node"
)

// WattsStrogatzGraph generates a Watts–Strogatz small-world graph.
// n = number of nodes
// k = each node is connected to k nearest neighbors in ring (must be even)
// beta = rewiring probability (0 = regular lattice, 1 = random graph)
func WattsStrogatzGraph(n, k int, beta float64, isUndirected bool) *graph.Graph {
	if n < 1 || k < 2 || k >= n || k%2 != 0 {
		return nil
	}

	ra := genRand()
	g := graph.New(isUndirected)

	// --- 1. Generate Nodes ---
	for i := 0; i < n; i++ {
		g.AddNode(node.ID(toString(i)))
	}

	// --- 2. Generate Regular Ring Lattice ---
	for i := 0; i < n; i++ {
		for j := 1; j <= k/2; j++ {
			neighbor := (i + j) % n
			g.AddEdge(node.ID(toString(i)), node.ID(toString(neighbor)))
		}
	}

	// --- 3. Rewiring Phase ---
	for i := 0; i < n; i++ {
		for j := 1; j <= k/2; j++ {
			neighbor := (i + j) % n
			if ra.Float64() < beta {
				// Remove existing edge
				g.RemoveEdge(node.ID(toString(i)), node.ID(toString(neighbor)))

				// Select a random other node (self-loop, duplicate prevention)
				for {
					newNeighbor := node.ID(toString(ra.Intn(n)))
					if newNeighbor != node.ID(toString(i)) && !g.HasEdge(node.ID(toString(i)), newNeighbor) {
						g.AddEdge(node.ID(toString(i)), newNeighbor)
						break
					}
				}
			}
		}
	}

	return g
}
