package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// DegreeCentrality calculates the degree centrality of a node in the graph.
func DegreeCentrality(g *graph.Graph, id node.ID, config *Config) float64 {
	if g == nil || !g.HasNode(id) {
		return 0.0
	}

	neighbors := g.GetNeighbors(id)

	return float64(len(neighbors))
}
