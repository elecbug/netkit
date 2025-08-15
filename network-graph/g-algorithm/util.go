package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
	"github.com/elecbug/go-dspkg/network-graph/path"
)

var cachedPaths map[string]map[node.ID]map[node.ID]path.Path = make(map[string]map[node.ID]map[node.ID]path.Path)

// AverageMetric computes the average value of a metric across all nodes in the graph.
func AverageMetric(g *graph.Graph, metricFunc func(*graph.Graph, node.ID, *Config) float64, config *Config) float64 {
	if len(g.GetNodes()) == 0 {
		return 0.0
	}

	total := 0.0
	count := 0

	for _, id := range g.GetNodes() {
		value := metricFunc(g, id, config)

		if value >= 0 {
			total += value
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return total / float64(count)
}
