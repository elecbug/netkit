package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// ClusteringCoefficient computes the clustering coefficient for a node in the graph.
func ClusteringCoefficient(g *graph.Graph, id node.ID, config *Config) float64 {
	if len(g.GetNodes()) == 0 {
		return 0.0
	}

	triangles := 0

	neighbors := g.GetNeighbors(id)
	numNeighbors := len(neighbors)

	if numNeighbors < 2 {
		return 0.0 // No triangles possible
	}

	for i := 0; i < numNeighbors; i++ {
		for j := i + 1; j < numNeighbors; j++ {
			if g.HasEdge(neighbors[i], neighbors[j]) {
				triangles++
			}
		}
	}

	return float64(triangles) / float64(numNeighbors*(numNeighbors-1)/2)
}
