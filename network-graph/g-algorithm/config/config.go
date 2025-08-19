package config

import "github.com/elecbug/go-dspkg/network-graph/node"

// Config holds the configuration settings for the graph algorithms.
type Config struct {
	Workers         int
	Closeness       *ClosenessCentralityConfig
	PageRank        *PageRankConfig
	EdgeBetweenness *EdgeBetweennessCentralityConfig
	Eigenvector     *EigenvectorCentralityConfig
	Degree          *DegreeCentralityConfig
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

// EdgeBetweennessCentralityConfig holds the configuration settings for the edge betweenness centrality algorithm.
type EdgeBetweennessCentralityConfig struct {
	Normalized bool
}

//
type EigenvectorCentralityConfig struct {
	MaxIter int
	Tol     float64
	Reverse bool
	NStart  *map[node.ID]float64 // initial vector; if nil, uniform distribution
}

//
type DegreeCentralityConfig struct {
	Mode string
}

// Default returns the default configuration for the graph algorithms.
func Default() *Config {
	return &Config{
		Workers:         16,
		Closeness:       &ClosenessCentralityConfig{WfImproved: true, Reverse: false},
		PageRank:        &PageRankConfig{Alpha: 0.85, MaxIter: 100, Tol: 1e-6},
		EdgeBetweenness: &EdgeBetweennessCentralityConfig{Normalized: true},
		Eigenvector:     &EigenvectorCentralityConfig{MaxIter: 100, Tol: 1e-6, Reverse: false, NStart: nil},
		Degree:          &DegreeCentralityConfig{Mode: "total"},
	}
}
