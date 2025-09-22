// Package config provides configuration settings for the graph algorithms.
package config

import (
	"runtime"
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
	Assortativity   *AssortativityCoefficientConfig
	Modularity      *ModularityConfig
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
		Assortativity:   &AssortativityCoefficientConfig{Mode: AssortativityProjected, IgnoreSelfLoops: true},
		Modularity:      &ModularityConfig{Partition: nil},
	}
}
