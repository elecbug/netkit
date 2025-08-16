package g_algorithm

import (
	"math"
	"slices"

	"github.com/elecbug/go-dspkg/network-graph/graph"
	"github.com/elecbug/go-dspkg/network-graph/node"
	"github.com/elecbug/go-dspkg/network-graph/path"
)

var cachedPaths map[string]path.GraphPaths = make(map[string]path.GraphPaths)

// ToGlobal computes the average value of a metric across all nodes in the graph.
func ToGlobal(g *graph.Graph, metricFunc func(*graph.Graph, node.ID, *Config) float64, config *Config) GlobalMetric {
	if len(g.GetNodes()) == 0 {
		return GlobalMetric{}
	}

	maps := make(map[node.ID]float64)
	total := []float64{}
	sum := 0.0
	count := 0

	for _, id := range g.GetNodes() {
		value := metricFunc(g, id, config)

		if value >= 0 {
			maps[id] = value
			total = append(total, value)
			sum += value
			count++
		}
	}

	if count == 0 {
		return GlobalMetric{}
	}

	return GlobalMetric{
		Average:           sum / float64(count),
		StandardDeviation: calculateStandardDeviation(total, sum/float64(count)),
		Middle:            total[len(total)/2],
		Max:               slices.Max(total),
		Min:               slices.Min(total),
		Values:            maps,
	}
}

func calculateStandardDeviation(total []float64, mean float64) float64 {
	if len(total) == 0 {
		return 0.0
	}

	var squaredDifferences []float64
	for _, value := range total {
		squaredDifferences = append(squaredDifferences, math.Pow(value-mean, 2))
	}

	return math.Sqrt(sum(squaredDifferences) / float64(len(squaredDifferences)))
}

func sum(values []float64) float64 {
	total := 0.0

	for _, v := range values {
		total += v
	}

	return total
}

type GlobalMetric struct {
	Average           float64
	StandardDeviation float64
	Middle            float64
	Max               float64
	Min               float64
	Values            map[node.ID]float64
}
