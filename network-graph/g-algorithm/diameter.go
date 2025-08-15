package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/path"
)

// Diameter computes the diameter of the graph using all-pairs shortest paths.
func Diameter(g *graph.Graph, workers int) path.Path {
	result := *path.NewPath()
	paths := AllShortestPaths(g, workers)

	for _, v := range paths {
		for _, p := range v {
			if p.GetDistance() > result.GetDistance() {
				result = p
			}
		}
	}

	return result
}
