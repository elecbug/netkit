package config

import "github.com/elecbug/go-dspkg/network-graph/node"

// Config holds the configuration settings for the graph algorithms.
type Config struct {
	Workers   int
	Closeness *ClosenessCentralityConfig
	PageRank  *PageRankConfig
}

// ClosenessCentralityConfig holds the configuration settings for the closeness centrality algorithm.
type ClosenessCentralityConfig struct {
	Reverse    bool
	WfImproved bool
}

// PageRankConfig holds the configuration settings for the PageRank algorithm.
type PageRankConfig struct {
	Alpha           float64              // damping, default 0.85
	MaxIter         int                  // default 100
	Tol             float64              // L1 error, default 1e-6
	Personalization *map[node.ID]float64 // p(u); if nil is uniform
	Dangling        *map[node.ID]float64 // d(u); if nil p(u)
	Reverse         bool
}
