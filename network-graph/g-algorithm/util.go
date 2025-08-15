package g_algorithm

import (
	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
)

// AverageMetric computes the average value of a metric across all nodes in the graph.
func AverageMetric(g *graph.Graph, metricFunc func(*graph.Graph, node.ID) float64) float64 {
	if len(g.GetNodes()) == 0 {
		return 0.0
	}

	total := 0.0
	count := 0

	for _, id := range g.GetNodes() {
		value := metricFunc(g, id)

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
