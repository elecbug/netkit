package standard_graph

import (
	"github.com/elecbug/netkit/graph"
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
		g.AddNode(graph.NodeID(toString(i)))
	}

	// --- 2. Generate Regular Ring Lattice ---
	for i := 0; i < n; i++ {
		for j := 1; j <= k/2; j++ {
			neighbor := (i + j) % n
			g.AddEdge(graph.NodeID(toString(i)), graph.NodeID(toString(neighbor)))
		}
	}

	// --- 3. Rewiring Phase ---
	for i := 0; i < n; i++ {
		for j := 1; j <= k/2; j++ {
			neighbor := (i + j) % n
			if ra.Float64() < beta {
				// Remove existing edge
				g.RemoveEdge(graph.NodeID(toString(i)), graph.NodeID(toString(neighbor)))

				// Select a random other node (self-loop, duplicate prevention)
				for {
					newNeighbor := graph.NodeID(toString(ra.Intn(n)))
					if newNeighbor != graph.NodeID(toString(i)) && !g.HasEdge(graph.NodeID(toString(i)), newNeighbor) {
						g.AddEdge(graph.NodeID(toString(i)), newNeighbor)
						break
					}
				}
			}
		}
	}

	return g
}
