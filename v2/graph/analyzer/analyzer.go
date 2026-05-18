package analyzer

import (
	"sync"

	"github.com/elecbug/netkit/v2/graph"
)

// Analyzer represents a graph analyzer that can be computed based on a given graph.
type Analyzer struct {
	baseGraph         *graph.Graph                                   // baseGraph is the original graph provided to the analyzer, used for reference and hashing.
	graphHash         string                                         // graphHash stores the hash of the base graph to detect changes and manage cache validity.
	allShortestPaths  map[graph.NodeID]map[graph.NodeID][]graph.Path // allShortestPaths caches the results of shortest path computations between node pairs.
	mu                sync.RWMutex                                   // mu protects access to the allShortestPaths cache to ensure thread safety during concurrent reads/writes.
	parallelCoreCount int                                            // parallelCoreCount determines how many CPU cores to utilize for parallel computations, if applicable.
	cfg               *Config                                        // cfg holds configuration options for the analyzer, such as worker counts and normalization settings.
}

// NewAnalyzer creates a new Analyzer instance based on the provided graph.
func NewAnalyzer(g *graph.Graph, parallelCoreCount int, cfg *Config) *Analyzer {
	return &Analyzer{
		baseGraph:         g,
		graphHash:         "",
		allShortestPaths:  make(map[graph.NodeID]map[graph.NodeID][]graph.Path),
		parallelCoreCount: parallelCoreCount,
		cfg:               cfg,
	}
}

// Graph returns the base graph associated with the analyzer.
func (a *Analyzer) Graph() *graph.Graph {
	return a.baseGraph
}
