package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// RichClubCoefficient computes the rich-club coefficient for a given node in the graph.
func RichClubCoefficient(g *graph.Graph, id node.ID, config *Config) float64 {
	// Compute the rich-club coefficient for the given node
	neighbors := g.GetNeighbors(id)
	if len(neighbors) == 0 {
		return 0.0
	}

	richClub := float64(len(neighbors)) / float64(len(g.GetNodes()))
	return richClub
}
