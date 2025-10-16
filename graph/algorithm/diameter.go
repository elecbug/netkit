package algorithm

import (
	"github.com/elecbug/netkit/graph"
)

// Diameter computes the diameter of the graph using all-pairs shortest paths.
func Diameter(g *graph.Graph, cfg *Config) int {
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
