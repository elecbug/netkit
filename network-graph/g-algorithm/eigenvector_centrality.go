package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// EigenvectorCentrality calculates the eigenvector centrality of a node in the graph.
func EigenvectorCentrality(g *graph.Graph, id node.ID, config *Config) float64 {
	if g == nil || !g.HasNode(id) {
		return 0.0
	}

	// Get the neighbors of the node
	neighbors := g.GetNeighbors(id)

	// Calculate the eigenvector centrality
	var centrality float64
	for _, neighbor := range neighbors {
		centrality += DegreeCentrality(g, neighbor, config)
	}

	return centrality
}
