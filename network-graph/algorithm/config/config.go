// Package config provides configuration settings for the graph algorithms.
package config

import (
	"runtime"

	"github.com/elecbug/netkit/network-graph/node"
)

// Config holds the configuration settings for the graph algorithms.
type Config struct {
	Workers         int
	Betweenness     *BetweennessCentralityConfig
	Closeness       *ClosenessCentralityConfig
	Degree          *DegreeCentralityConfig
	EdgeBetweenness *EdgeBetweennessCentralityConfig
	Eigenvector     *EigenvectorCentralityConfig
	PageRank        *PageRankConfig
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

// BetweennessCentralityConfig holds the configuration settings for the edge betweenness centrality algorithm.
type BetweennessCentralityConfig struct {
	Normalized bool
}

// EdgeBetweennessCentralityConfig holds the configuration settings for the edge betweenness centrality algorithm.
type EdgeBetweennessCentralityConfig struct {
	Normalized bool
}

// EigenvectorCentralityConfig holds the configuration settings for the eigenvector centrality algorithm.
type EigenvectorCentralityConfig struct {
	MaxIter int
	Tol     float64
	Reverse bool
	NStart  *map[node.ID]float64 // initial vector; if nil, uniform distribution
}

// DegreeCentralityConfig holds the configuration settings for the degree centrality algorithm.
type DegreeCentralityConfig struct {
	Mode string
}

// Default returns the default configuration for the graph algorithms.
func Default() *Config {
	return &Config{
		Workers:         runtime.NumCPU(),
		Closeness:       &ClosenessCentralityConfig{WfImproved: true, Reverse: false},
		PageRank:        &PageRankConfig{Alpha: 0.85, MaxIter: 100, Tol: 1e-6, Personalization: nil, Dangling: nil, Reverse: false},
		Betweenness:     &BetweennessCentralityConfig{Normalized: true},
		EdgeBetweenness: &EdgeBetweennessCentralityConfig{Normalized: true},
		Eigenvector:     &EigenvectorCentralityConfig{MaxIter: 100, Tol: 1e-6, Reverse: false, NStart: nil},
		Degree:          &DegreeCentralityConfig{Mode: "total"},
	}
}
