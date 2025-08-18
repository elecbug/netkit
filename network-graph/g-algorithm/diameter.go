package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/g-algorithm/config"
	"github.com/elecbug/go-dspkg/network-graph/graph"
)

// Diameter computes the diameter of the graph using all-pairs shortest paths.
func Diameter(g *graph.Graph, cfg *config.Config) int {
	result := 0
	paths := AllShortestPaths(g, cfg)

	for _, v := range paths {
		for _, ps := range v {
			if len(ps) == 0 {
				continue
			}

			if ps[0].Distance() > result {
				result = ps[0].Distance()
			}
		}
	}

	return result
}
