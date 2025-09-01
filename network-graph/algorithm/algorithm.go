// Package algorithm provides graph algorithms for network analysis.
package algorithm

import (
	"github.com/elecbug/netkit/network-graph/algorithm/config"
	"github.com/elecbug/netkit/network-graph/graph"
)

// MetricType represents the type of metric to be calculated.
type MetricType int

// Metric types for graph analysis.
const (
	BETWEENNESS_CENTRALITY MetricType = iota
	CLOSENESS_CENTRALITY
	CLUSTERING_COEFFICIENT
	DEGREE_CENTRALITY
	DIAMETER
	EDGE_BETWEENNESS_CENTRALITY
	EIGENVECTOR_CENTRALITY
	PAGE_RANK
	SHORTEST_PATHS
)

// Metric calculates the specified metric for the given graph.
func Metric(g *graph.Graph, cfg *config.Config, metricType MetricType) any {
	switch metricType {
	case BETWEENNESS_CENTRALITY:
		return BetweennessCentrality(g, cfg)
	case CLOSENESS_CENTRALITY:
		return ClosenessCentrality(g, cfg)
	case CLUSTERING_COEFFICIENT:
		return ClusteringCoefficient(g, cfg)
	case DEGREE_CENTRALITY:
		return DegreeCentrality(g, cfg)
	case DIAMETER:
		return Diameter(g, cfg)
	case EDGE_BETWEENNESS_CENTRALITY:
		return EdgeBetweennessCentrality(g, cfg)
	case EIGENVECTOR_CENTRALITY:
		return EigenvectorCentrality(g, cfg)
	case PAGE_RANK:
		return PageRank(g, cfg)
	case SHORTEST_PATHS:
		return AllShortestPaths(g, cfg)
	default:
		return nil
	}
}
