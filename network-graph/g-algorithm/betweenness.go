package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// BetweennessCentrality computes the betweenness centrality for a node in the graph.
func BetweennessCentrality(g *graph.Graph, id node.ID) float64 {
	if len(g.GetNodes()) == 0 {
		return 0.0
	}

	betweenness := 0.0
	nodes := g.GetNodes()

	for _, s := range nodes {
		if s == id {
			continue // Skip the node itself
		}

		for _, e := range nodes {
			if e == id || e == s {
				continue // Skip if it's the same node or the source node
			}

			paths := AllShortestPaths(g, 1) // Assuming single-threaded for simplicity
			p := paths[s][e]
			for _, node := range p.GetNodes() {
				if node == id {
					betweenness++
				}
			}
		}
	}

	return betweenness / float64(len(nodes)-1)
}
