package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

func LocalEfficiency(g *graph.Graph, id node.ID, config *Config) float64 {
	neighbors := g.GetNeighbors(id)
	if len(neighbors) < 2 {
		return 0.0
	}

	subgraph := graph.New(g.IsBidirectional())

	for _, n := range neighbors {
		subgraph.AddNode(n)
	}

	for i := 0; i < len(neighbors); i++ {
		for j := i + 1; j < len(neighbors); j++ {
			if g.HasEdge(neighbors[i], neighbors[j]) {
				subgraph.AddEdge(neighbors[i], neighbors[j])
			}
		}
	}

	actualEdges := float64(subgraph.EdgeCount())
	possibleEdges := float64(len(neighbors) * (len(neighbors) - 1) / 2)

	if possibleEdges == 0 {
		return 0.0
	}

	return actualEdges / possibleEdges
}
