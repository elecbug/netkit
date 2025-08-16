package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
)

// Diameter computes the diameter of the graph using all-pairs shortest paths.
func Diameter(g *graph.Graph, config *Config) int {
	result := 0
	paths := AllShortestPaths(g, config)

	for _, v := range paths {
		for _, ps := range v {
			if ps[0].GetDistance() > result {
				result = ps[0].GetDistance()
			}
		}
	}

	return result
}
