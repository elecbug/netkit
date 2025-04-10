package graph_algorithm

import (
	"github.com/elecbug/go-dspkg/graph"
)

// `Unit` represents a computation unit for graph algorithms.
// It stores shortest paths within the graph and performs computations.
type Unit struct {
	shortestPaths []Path       // Stores the shortest paths for the graph, sorted by distance in ascending order.
	graph         *graph.Graph // A reference to the graph associated with this computation unit.
	updateVersion int          // Update information for shortest paths
}

// `ParallelUnit` extends Unit for parallel computation of graph algorithms.
// It supports dividing tasks across multiple CPU cores for better performance.
type ParallelUnit struct {
	Unit         // Embeds the Unit structure for shared functionality.
	maxCore uint // Maximum number of CPU cores to use for parallel computation.
}

// `NewUnit` creates and initializes a new `Unit` instance.
func NewUnit(g *graph.Graph) *Unit {
	return &Unit{
		shortestPaths: make([]Path, g.EdgeCount()), // Initialize with an empty slice of paths.
		graph:         g,                           // Associate the graph with this Unit.
		updateVersion: -1,
	}
}

// `NewParallelUnit` creates and initializes a new `ParallelUnit` instance.
func NewParallelUnit(g *graph.Graph, core uint) *ParallelUnit {
	return &ParallelUnit{
		Unit:    *NewUnit(g), // Initialize the embedded Unit structure.
		maxCore: core,        // Set the maximum number of cores for parallel processing.
	}
}
